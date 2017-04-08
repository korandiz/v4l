// Package v4l, a facade to the Video4Linux video capture interface
// Copyright (C) 2016 Zoltán Korándi <korandi.z@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Command viewcam displays the video captured from a V4L device in a GTK-based
// GUI.
//
// Usage:
//   viewcam
package main

import (
	"errors"
	"fmt"
	"image"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/yuyv"
)

func main() {
	var (
		dp DevicePickerWindow
		vc VideoCaptureWindow
		cp ControlPanelWindow
	)

	dp.OnClose = func() {
		if dp.Video {
			vc.DevicePath = dp.DevicePath
			vc.DeviceConfig = dp.DeviceConfig
			vc.Open()
		}
		if dp.ControlPanel {
			cp.DevicePath = dp.DevicePath
			cp.Open()
		}
		if !dp.Video && !dp.ControlPanel {
			gtk.MainQuit()
		}
	}

	vc.OnClose = func() {
		cp.Close()
		dp.Open()
	}

	cp.OnClose = func() {
		vc.Close()
		dp.Open()
	}

	gtk.Init(nil)
	dp.Open()
	gtk.Main()
}

type DevicePickerWindow struct {
	DevicePath   string
	DeviceConfig v4l.DeviceConfig
	Video        bool
	ControlPanel bool
	OnClose      func()

	open    bool
	devices []v4l.DeviceInfo
	configs []v4l.DeviceConfig
	ticker  *ticker

	window       *gtk.Window
	deviceList   *gtk.ListBox
	configList   *gtk.ListBox
	pathLabel    *gtk.Label
	nameLabel    *gtk.Label
	businfoLabel *gtk.Label
	driverLabel  *gtk.Label
	cameraLabel  *gtk.Label
	videoButton  *gtk.RadioButton
	ctrlButton   *gtk.RadioButton
	bothButton   *gtk.RadioButton
}

func (w *DevicePickerWindow) Open() {
	if w.open {
		return
	}
	w.DevicePath = ""
	w.DeviceConfig = v4l.DeviceConfig{}
	w.Video = false
	w.ControlPanel = false
	w.initWidgets()
	w.open = true
	w.updateDeviceList()
	w.ticker = newTicker(500, w.updateDeviceList)
}

func (w *DevicePickerWindow) Close() {
	if !w.open {
		return
	}
	w.window.Destroy()
}

func (w *DevicePickerWindow) initWidgets() {
	var err error

	w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fatal(err)
	w.window.SetTitle("Video Capture Demo")
	w.window.SetIconName("camera-web")
	w.window.SetPosition(gtk.WIN_POS_CENTER)
	w.window.SetSizeRequest(640, 320)
	w.window.SetBorderWidth(6)
	_, err = w.window.Connect("destroy", w.windowDestroyed)
	fatal(err)

	grid1, err := gtk.GridNew()
	fatal(err)
	grid1.SetColumnHomogeneous(true)
	grid1.SetColumnSpacing(6)
	grid1.SetRowSpacing(6)
	w.window.Add(grid1)

	frame, err := gtk.FrameNew("Devices:")
	fatal(err)
	grid1.Attach(frame, 0, 0, 1, 1)

	scrolled, err := gtk.ScrolledWindowNew(nil, nil)
	fatal(err)
	scrolled.SetHExpand(true)
	scrolled.SetVExpand(true)
	frame.Add(scrolled)

	w.deviceList, err = gtk.ListBoxNew()
	fatal(err)
	_, err = w.deviceList.Connect("selected-rows-changed", w.deviceSelected)
	fatal(err)
	scrolled.Add(w.deviceList)

	frame, err = gtk.FrameNew("Configurations:")
	fatal(err)
	grid1.Attach(frame, 1, 0, 1, 1)

	scrolled, err = gtk.ScrolledWindowNew(nil, nil)
	fatal(err)
	scrolled.SetHExpand(true)
	scrolled.SetVExpand(true)
	frame.Add(scrolled)

	w.configList, err = gtk.ListBoxNew()
	fatal(err)
	_, err = w.configList.Connect("selected-rows-changed", w.configSelected)
	fatal(err)
	scrolled.Add(w.configList)

	grid2, err := gtk.GridNew()
	fatal(err)
	grid2.SetColumnSpacing(12)
	grid2.SetRowSpacing(6)
	grid1.Attach(grid2, 0, 1, 1, 1)

	w.pathLabel = labelNew("")
	w.nameLabel = labelNew("")
	w.businfoLabel = labelNew("")
	w.driverLabel = labelNew("")
	w.cameraLabel = labelNew("")
	grid2.Attach(labelNew("Path:"), 0, 0, 1, 1)
	grid2.Attach(labelNew("Name:"), 0, 1, 1, 1)
	grid2.Attach(labelNew("BusInfo:"), 0, 2, 1, 1)
	grid2.Attach(labelNew("Driver:"), 0, 3, 1, 1)
	grid2.Attach(labelNew("Camera:"), 0, 4, 1, 1)
	grid2.Attach(w.pathLabel, 1, 0, 1, 1)
	grid2.Attach(w.nameLabel, 1, 1, 1, 1)
	grid2.Attach(w.businfoLabel, 1, 2, 1, 1)
	grid2.Attach(w.driverLabel, 1, 3, 1, 1)
	grid2.Attach(w.cameraLabel, 1, 4, 1, 1)
	w.clearLabels()

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	fatal(err)
	grid1.Attach(box, 1, 1, 1, 1)

	w.videoButton, err = gtk.RadioButtonNewWithLabel(nil, "Video")
	fatal(err)
	box.Add(w.videoButton)

	w.ctrlButton, err = gtk.RadioButtonNewWithLabelFromWidget(w.videoButton,
		"Control Panel")
	fatal(err)
	box.Add(w.ctrlButton)

	w.bothButton, err = gtk.RadioButtonNewWithLabelFromWidget(w.videoButton,
		"Both")
	fatal(err)
	box.Add(w.bothButton)

	button, err := gtk.ButtonNewWithLabel("Start")
	fatal(err)
	_, err = button.Connect("clicked", w.startClicked)
	fatal(err)
	box.Add(button)

	w.window.ShowAll()
}

func (w *DevicePickerWindow) windowDestroyed() {
	w.ticker.cancel()
	*w = DevicePickerWindow{
		DevicePath:   w.DevicePath,
		DeviceConfig: w.DeviceConfig,
		Video:        w.Video,
		ControlPanel: w.ControlPanel,
		OnClose:      w.OnClose,
	}
	w.OnClose()
}

func (w *DevicePickerWindow) updateDeviceList() {
	devices := v4l.FindDevices()
	if deviceListsEqual(devices, w.devices) {
		return
	}
	oldIndex := -1
	newIndex := -1
	if row := w.deviceList.GetSelectedRow(); row != nil {
		oldIndex = row.GetIndex()
		for i := range devices {
			if devices[i] == w.devices[oldIndex] {
				newIndex = i
			}
		}
	}
	for i := range w.devices {
		if i < oldIndex || newIndex == -1 {
			w.deviceList.Remove(w.deviceList.GetRowAtIndex(0))
		} else if i > oldIndex {
			w.deviceList.Remove(w.deviceList.GetRowAtIndex(1))
		}
	}
	for i := range devices {
		if i < newIndex {
			w.deviceList.Insert(labelNew(devices[i].Path), i)
		} else if i > newIndex {
			w.deviceList.Add(labelNew(devices[i].Path))
		}
	}
	w.devices = devices
	w.deviceList.ShowAll()
}

func (w *DevicePickerWindow) deviceSelected() {
	if !w.open {
		return
	}
	w.clearConfigList()
	row := w.deviceList.GetSelectedRow()
	if row == nil {
		w.clearLabels()
		w.DevicePath = ""
		return
	}
	err := w.fillConfigList(w.devices[row.GetIndex()])
	if err != nil {
		w.modalWarning(err.Error())
		w.deviceList.SelectRow(nil)
		return
	}
	device := w.devices[row.GetIndex()]
	w.setLabels(device)
	w.DevicePath = device.Path
}

func (w *DevicePickerWindow) clearConfigList() {
	w.configs = nil
	w.DeviceConfig = v4l.DeviceConfig{}
	for {
		row := w.configList.GetRowAtIndex(0)
		if row == nil {
			break
		}
		w.configList.Remove(row)
	}
}

func (w *DevicePickerWindow) clearLabels() {
	w.pathLabel.SetLabel("\u2014")
	w.nameLabel.SetLabel("\u2014")
	w.businfoLabel.SetLabel("\u2014")
	w.driverLabel.SetLabel("\u2014")
	w.cameraLabel.SetLabel("\u2014")
}

func (w *DevicePickerWindow) setLabels(info v4l.DeviceInfo) {
	dn := info.DriverName
	dv := info.DriverVersion
	c := "No"
	if info.Camera {
		c = "Yes"
	}
	w.pathLabel.SetLabel(info.Path)
	w.nameLabel.SetLabel(info.DeviceName)
	w.businfoLabel.SetLabel(info.BusInfo)
	w.driverLabel.SetLabel(fmt.Sprintf("%s %v.%v.%v", dn, dv[0], dv[1], dv[2]))
	w.cameraLabel.SetLabel(c)
}

func (w *DevicePickerWindow) fillConfigList(info v4l.DeviceInfo) error {
	dev, err := v4l.Open(info.Path)
	if err != nil {
		return errors.New("Open: " + err.Error())
	}
	defer dev.Close()
	cfgs, err := dev.ListConfigs()
	if err != nil {
		return errors.New("ListConfigs: " + err.Error())
	}
	for i := range cfgs {
		if cfgs[i].Format != yuyv.FourCC {
			continue
		}
		w.configList.Add(labelNew(cfg2str(cfgs[i])))
		w.configs = append(w.configs, cfgs[i])
	}
	w.configList.ShowAll()
	return nil
}

func (w *DevicePickerWindow) configSelected() {
	if !w.open {
		return
	}
	row := w.configList.GetSelectedRow()
	if row == nil {
		w.DeviceConfig = v4l.DeviceConfig{}
		return
	}
	w.DeviceConfig = w.configs[row.GetIndex()]
}

func (w *DevicePickerWindow) startClicked() {
	video := w.videoButton.GetActive() || w.bothButton.GetActive()
	controlPanel := w.ctrlButton.GetActive() || w.bothButton.GetActive()
	if row := w.deviceList.GetSelectedRow(); row == nil {
		w.modalWarning("Please select a device.")
		return
	}
	if row := w.configList.GetSelectedRow(); row == nil && video {
		w.modalWarning("Please select a configuration.")
		return
	}
	w.Video = video
	w.ControlPanel = controlPanel
	w.window.Destroy()
}

func (w *DevicePickerWindow) modalWarning(msg string) {
	dlg := gtk.MessageDialogNew(w.window,
		gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, msg)
	dlg.SetTitle("Error")
	dlg.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	_, err := dlg.Connect("response", dlg.Destroy)
	fatal(err)
	dlg.ShowAll()
}

type VideoCaptureWindow struct {
	DevicePath   string
	DeviceConfig v4l.DeviceConfig
	OnClose      func()

	open    bool
	capture *capture

	window *gtk.Window
	pixbuf *gdk.Pixbuf
	image  *gtk.Image
}

func (w *VideoCaptureWindow) Open() {
	if w.open {
		return
	}
	w.initWidgets()
	w.open = true
	w.capture = newCapture(w)
}

func (w *VideoCaptureWindow) Close() {
	if !w.open {
		return
	}
	w.window.Destroy()
}

func (w *VideoCaptureWindow) initWidgets() {
	var err error

	w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fatal(err)
	w.window.SetTitle("Video Capture Demo")
	w.window.SetIconName("camera-web")
	w.window.SetPosition(gtk.WIN_POS_CENTER)
	w.window.SetResizable(false)
	_, err = w.window.Connect("destroy", w.windowDestroyed)
	fatal(err)

	w.pixbuf, err = gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8,
		w.DeviceConfig.Width, w.DeviceConfig.Height)
	fatal(err)
	for p := w.pixbuf.GetPixels(); len(p) > 0; p = p[1:] {
		p[0] = 0
	}

	w.image, err = gtk.ImageNew()
	fatal(err)
	w.image.SetFromPixbuf(w.pixbuf)
	w.window.Add(w.image)

	w.window.ShowAll()
}

func (w *VideoCaptureWindow) windowDestroyed() {
	w.capture.cancel()
	*w = VideoCaptureWindow{
		DevicePath:   w.DevicePath,
		DeviceConfig: w.DeviceConfig,
		OnClose:      w.OnClose,
	}
	w.OnClose()
}

func (w *VideoCaptureWindow) modalError(prefix string, err error) {
	msg := prefix + ": " + err.Error()
	dlg := gtk.MessageDialogNew(w.window,
		gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	dlg.SetTitle("Error")
	dlg.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	_, err1 := dlg.Connect("response", dlg.Destroy)
	fatal(err1)
	_, err1 = dlg.Connect("destroy", w.window.Destroy)
	fatal(err1)
	dlg.ShowAll()
}

type capture struct {
	cancelled bool
	w         *VideoCaptureWindow
	yuyv      *yuyv.Image
}

func newCapture(w *VideoCaptureWindow) *capture {
	c := &capture{w: w}
	go c.captureFrames(w.DevicePath, w.DeviceConfig)
	return c
}

func (c *capture) cancel() {
	c.cancelled = true
}

func (c *capture) captureFrames(path string, config v4l.DeviceConfig) {
	device, err := v4l.Open(path)
	if err != nil {
		c.modalError("Open", err)
		return
	}
	defer device.Close()
	if err := device.SetConfig(config); err != nil {
		c.modalError("SetConfig", err)
		return
	}
	if err := device.TurnOn(); err != nil {
		c.modalError("TurnOn", err)
		return
	}
	config, err = device.GetConfig()
	if err != nil {
		c.modalError("GetConfig", err)
		return
	}
	if config.Format != yuyv.FourCC {
		c.modalError("SetConfig", errors.New("failed to set YUYV format"))
		return
	}
	binfo, err := device.BufferInfo()
	if err != nil {
		c.modalError("BufferInfo", err)
		return
	}
	c.yuyv = &yuyv.Image{
		Pix:    make([]byte, binfo.BufferSize),
		Stride: binfo.ImageStride,
		Rect:   image.Rect(0, 0, config.Width, config.Height),
	}
	for {
		buffer, err := device.Capture()
		if err != nil {
			c.modalError("Capture", err)
			return
		}
		buffer.ReadAt(c.yuyv.Pix, 0)
		ok := c.updateImage()
		if !ok {
			return
		}
	}
}

func (c *capture) modalError(prefix string, err error) {
	ch := make(chan struct{})
	_, err1 := glib.IdleAdd(func() {
		if !c.cancelled {
			c.w.modalError(prefix, err)
		}
		ch <- struct{}{}
	})
	fatal(err1)
	<-ch
}

func (c *capture) updateImage() bool {
	ch := make(chan bool)
	_, err := glib.IdleAdd(func() {
		if c.cancelled {
			ch <- false
			return
		}
		w := c.w
		rgba := &image.RGBA{
			Pix:    w.pixbuf.GetPixels(),
			Stride: w.pixbuf.GetRowstride(),
			Rect:   image.Rect(0, 0, w.pixbuf.GetWidth(), w.pixbuf.GetHeight()),
		}
		yuyv.ToRGBA(rgba, rgba.Rect, c.yuyv, c.yuyv.Rect.Min)
		w.image.SetFromPixbuf(c.w.pixbuf)
		ch <- true
	})
	fatal(err)
	return <-ch
}

type ControlPanelWindow struct {
	DevicePath string
	OnClose    func()

	open     bool
	device   *v4l.Device
	controls []v4l.ControlInfo
	updaters []updater
	ticker   *ticker

	window *gtk.Window
}

type updater func(int32)

func (w *ControlPanelWindow) Open() {
	if w.open {
		return
	}
	w.initWidgets()
	w.open = true
	w.initControls()
	w.updateControls()
	w.ticker = newTicker(100, w.updateControls)
}

func (w *ControlPanelWindow) Close() {
	if !w.open {
		return
	}
	w.window.Destroy()
}

func (w *ControlPanelWindow) initWidgets() {
	var err error
	w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fatal(err)
	w.window.SetTitle("Control Panel")
	w.window.SetIconName("camera-web")
	w.window.SetPosition(gtk.WIN_POS_CENTER)
	w.window.SetBorderWidth(6)
	w.window.SetSizeRequest(640, 0)
	_, err = w.window.Connect("destroy", w.windowDestroyed)
	fatal(err)
	w.window.ShowAll()
}

func (w *ControlPanelWindow) initControls() {
	var err error
	w.device, err = v4l.Open(w.DevicePath)
	if err != nil {
		w.modalError("Open", err)
		return
	}
	w.controls, err = w.device.ListControls()
	if err != nil {
		w.modalError("ListControls", err)
		return
	}
	grid, err := gtk.GridNew()
	fatal(err)
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(6)
	grid.SetRowHomogeneous(true)
	w.window.Add(grid)
	for i, ctrl := range w.controls {
		frame, err := gtk.FrameNew("")
		fatal(err)
		if ctrl.Type != "button" {
			grid.Attach(labelNew(ctrl.Name), 0, i, 1, 1)
		}
		if ctrl.Type != "int" {
			frame.SetHAlign(gtk.ALIGN_START)
			frame.SetVAlign(gtk.ALIGN_CENTER)
		}
		frame.SetShadowType(gtk.SHADOW_NONE)
		grid.Attach(frame, 1, i, 1, 1)
		var upd updater
		switch ctrl.Type {
		case "int":
			upd = w.createIntControl(&frame.Container, ctrl)
		case "bool":
			upd = w.createBoolControl(&frame.Container, ctrl)
		case "enum":
			upd = w.createEnumControl(&frame.Container, ctrl)
		case "button":
			upd = w.createButtonControl(&frame.Container, ctrl)
		}
		w.updaters = append(w.updaters, upd)
	}
	grid.ShowAll()
}

func (w *ControlPanelWindow) createIntControl(parent *gtk.Container, ctrl v4l.ControlInfo) updater {
	var err error
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	fatal(err)
	parent.Add(box)
	scale, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL,
		float64(ctrl.Min), float64(ctrl.Max), float64(ctrl.Step))
	fatal(err)
	scale.SetHExpand(true)
	scale.SetProperty("draw-value", false)
	box.Add(scale)
	frame, err := gtk.FrameNew("")
	fatal(err)
	frame.SetShadowType(gtk.SHADOW_NONE)
	frame.SetSizeRequest(50, 0)
	box.Add(frame)
	label, err := gtk.LabelNew("")
	fatal(err)
	label.SetHAlign(gtk.ALIGN_END)
	frame.Add(label)
	_, err = scale.Connect("value-changed", func() {
		value := int32(scale.GetValue() + 0.5)
		w.device.SetControl(ctrl.CID, value)
		w.updateControls()
	})
	fatal(err)
	return func(value int32) {
		scale.SetValue(float64(value))
		label.SetLabel(strconv.Itoa(int(value)))
	}
}

func (w *ControlPanelWindow) createBoolControl(parent *gtk.Container, ctrl v4l.ControlInfo) updater {
	checkButton, err := gtk.CheckButtonNew()
	fatal(err)
	checkButton.SetHExpand(false)
	checkButton.SetVExpand(false)
	parent.Add(checkButton)
	_, err = checkButton.Connect("toggled", func() {
		value := int32(0)
		if checkButton.GetActive() {
			value = 1
		}
		w.device.SetControl(ctrl.CID, value)
		w.updateControls()
	})
	fatal(err)
	return func(value int32) {
		active := false
		if value != 0 {
			active = true
		}
		checkButton.SetActive(active)
	}
}

func (w *ControlPanelWindow) createEnumControl(parent *gtk.Container, ctrl v4l.ControlInfo) updater {
	comboBox, err := gtk.ComboBoxTextNew()
	fatal(err)
	parent.Add(comboBox)
	for _, opt := range ctrl.Options {
		comboBox.Append("", opt.Name)
	}
	_, err = comboBox.Connect("changed", func() {
		value := ctrl.Options[comboBox.GetActive()].Value
		w.device.SetControl(ctrl.CID, value)
		w.updateControls()
	})
	fatal(err)
	return func(value int32) {
		for i, opt := range ctrl.Options {
			if opt.Value == value {
				comboBox.SetActive(i)
				break
			}
		}
	}
}

func (w *ControlPanelWindow) createButtonControl(parent *gtk.Container, ctrl v4l.ControlInfo) updater {
	button, err := gtk.ButtonNewWithLabel(ctrl.Name)
	fatal(err)
	parent.Add(button)
	_, err = button.Connect("clicked", func() {
		w.device.SetControl(ctrl.CID, 0)
		w.updateControls()
	})
	fatal(err)
	return nil
}

func (w *ControlPanelWindow) windowDestroyed() {
	w.ticker.cancel()
	if w.device != nil {
		w.device.Close()
	}
	*w = ControlPanelWindow{
		DevicePath: w.DevicePath,
		OnClose:    w.OnClose,
	}
	w.OnClose()
}

func (w *ControlPanelWindow) updateControls() {
	for i, ctrl := range w.controls {
		if ctrl.Type == "button" {
			continue
		}
		value, err := w.device.GetControl(ctrl.CID)
		if err != nil {
			continue
		}
		if w.updaters[i] != nil {
			w.updaters[i](value)
		}
	}
}

func (w *ControlPanelWindow) modalError(prefix string, err error) {
	msg := prefix + ": " + err.Error()
	dlg := gtk.MessageDialogNew(w.window,
		gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	dlg.SetTitle("Error")
	dlg.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	_, err1 := dlg.Connect("response", dlg.Destroy)
	fatal(err1)
	_, err1 = dlg.Connect("destroy", w.window.Destroy)
	fatal(err1)
	dlg.ShowAll()
}

type ticker struct {
	interval  uint
	callback  func()
	cancelled bool
}

func newTicker(interval uint, callback func()) *ticker {
	t := &ticker{
		interval: interval,
		callback: callback,
	}
	_, err := glib.TimeoutAdd(interval, t.handler)
	fatal(err)
	return t
}

func (t *ticker) cancel() {
	t.cancelled = true
}

func (t *ticker) handler() {
	if t.cancelled {
		return
	}
	t.callback()
	_, err := glib.TimeoutAdd(t.interval, t.handler)
	fatal(err)

}

func deviceListsEqual(d1, d2 []v4l.DeviceInfo) bool {
	if len(d1) != len(d2) {
		return false
	}
	for i := range d1 {
		if d1[i] != d2[i] {
			return false
		}
	}
	return true
}

func cfg2str(cfg v4l.DeviceConfig) string {
	w := cfg.Width
	h := cfg.Height
	f := cfg.FPS
	return fmt.Sprintf("%dx%d @ %.4g FPS", w, h, float64(f.N)/float64(f.D))
}

func labelNew(s string) *gtk.Label {
	l, err := gtk.LabelNew(s)
	fatal(err)
	l.SetHAlign(gtk.ALIGN_START)
	return l
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
