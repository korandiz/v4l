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

//go:build linux
// +build linux

// Package v4l is a facade to the Video4Linux video capture interface.
package v4l

import "syscall"

// Control IDs. Devices may have other controls than these, including custom
// (driver specific) ones.
const (
	// Integer controls
	CtrlBrightness            = 0x00980900
	CtrlContrast              = 0x00980901
	CtrlSaturation            = 0x00980902
	CtrlHue                   = 0x00980903
	CtrlGamma                 = 0x00980910
	CtrlExposure              = 0x00980911
	CtrlGain                  = 0x00980913
	CtrlWhiteBalance          = 0x0098091a
	CtrlSharpness             = 0x0098091b
	CtrlBacklightCompensation = 0x0098091c

	// Boolean controls
	CtrlHFlip            = 0x00980914
	CtrlVFlip            = 0x00980915
	CtrlAutoWhiteBalance = 0x0098090c
	CtrlAutoGain         = 0x00980912
	CtrlAutoHue          = 0x00980919
	CtrlAutoBrightness   = 0x00980920

	// Enums
	CtrlPowerLineFreq     = 0x00980918
	PowerLineFreqDisabled = 0
	PowerLineFreq50Hz     = 1
	PowerLineFreq60Hz     = 2
	PowerLineFreqAuto     = 3

	// Buttons
	CtrlDoWhiteBalance = 0x0098090d
)

// A Device represents a V4L capture device.
type Device struct {
	*device
}

// device is the real representation of Device. The extra level of indirection
// is there to prevent clients from tampering with the file descriptor and the
// buffers.
type device struct {
	path      string
	fd        int
	buffers   [][]byte
	bufIndex  uint32
	nCaptures uint64
}

// noBuffer is the value assinged to device.bufIndex when none of the buffers
// contain valid image data.
const noBuffer = ^uint32(0)

// A DeviceInfo provides information about a capture device.
type DeviceInfo struct {
	// Path is the device path. (e.g. /dev/video0)
	Path string

	// DeviceName is the name of the device. (e.g. "Yoyodyne TV/FM")
	DeviceName string

	// BusInfo is the location of the device in the system.
	// (e.g. "PCI:0000:05:06.0")
	BusInfo string

	// DriverName is the name of the driver. (e.g. "bttv")
	DriverName string

	// DriverVersion contains the three components of the driver's version
	// number. (e.g. [3]int{1, 2, 3} for version 1.2.3)
	DriverVersion [3]int

	// Camera tells if the device is a camera. If false, then this is some other
	// kind of capture device, e.g. an analog TV tuner.
	Camera bool
}

// A DeviceConfig encapsulates the configuration of a capture device.
type DeviceConfig struct {
	// Format is the four-character code (FourCC) of the pixel format, with the
	// first character at the lowest byte.
	Format uint32

	// Width and Height specify the image dimensions.
	Width  int
	Height int

	// FPS specifies the frame rate.
	FPS Frac
}

// A BufferInfo provides information about how image data is laid out in a
// buffer.
type BufferInfo struct {
	// BufferSize is the number of bytes required to hold an image. For variable
	// length compressed formats, it's the maximum size an image may take up.
	BufferSize int

	// ImageStride is the distance in bytes between the leftmost pixels of
	// adjacent lines.
	ImageStride int
}

// A ControlInfo provides information about a control.
type ControlInfo struct {
	// CID is the identifier of the control. (e.g. 0x00980900)
	CID uint32

	// Name is the name of the control. (e.g. "Brightness")
	Name string

	// Type is the type of the control, one of "int", "bool", "enum",
	// or "button".
	//   - The valid values of an integer control are determined by Min, Max,
	//     and Step.
	//   - A boolean control can only have the values 0 and 1, where 0 means
	//     "disabled" and 1 means "enabled".
	//   - Enums can only take values from a predefined set. (see Options)
	//   - Buttons perform some action when pushed, and they don't have a value.
	//     Reading the value of a button fails, while setting it to any value is
	//     interpreted as a push.
	Type string

	// Min and Max specify the range of values the control can take, and Step is
	// the smallest change actually affecting the hardware. They are only
	// meaningful for integer type controls.
	Min  int32
	Max  int32
	Step int32

	// Default is the default value of the control.
	Default int32

	// Options is the list of valid values of an enum type control. For other
	// types it's nil.
	Options []struct {
		Value int32
		Name  string
	}
}

// errBadControl is returned by Device.controlInfo when the control is disabled
// or of unsupported type.
const errBadControl = Error("control disabled or of unsupported type")

// Open opens the capture device named by path. If the file is not a capture
// device, it fails with ErrWrongDevice.
func Open(path string) (*Device, error) {
	// Open the file.
	fd, err := syscall.Open(path, syscall.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}

	// Check if it's a V4L device.
	var stat syscall.Stat_t
	if err := syscall.Fstat(fd, &stat); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	if stat.Mode&syscall.S_IFCHR == 0 || stat.Rdev>>8 != 81 {
		syscall.Close(fd)
		return nil, ErrWrongDevice
	}

	// Check if it's a capture device.
	var c v4l_capability
	if err := ioctl_querycap(fd, &c); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	caps := c.capabilities
	if caps&v4l_capDeviceCaps != 0 {
		caps = c.deviceCaps
	}
	if caps&v4l_capVideoCapture == 0 {
		syscall.Close(fd)
		return nil, ErrWrongDevice
	}

	d := device{path: path, fd: fd, bufIndex: noBuffer}
	return &Device{&d}, nil
}

// Close closes the device, freeing all native resources associated with it. It
// stops any capture session in progress, and it may also render the contents of
// previously captured buffers unavailable.
func (d *Device) Close() {
	d.TurnOff()
	syscall.Close(d.fd)
	d.fd = -1
}

// DeviceInfo returns information about the device.
func (d *Device) DeviceInfo() (DeviceInfo, error) {
	// Query capabilities.
	var c v4l_capability
	if err := ioctl_querycap(d.fd, &c); err != nil {
		return DeviceInfo{}, err
	}

	// Check if the device is a camera. Cameras do not enumerate any video
	// standards.
	s := v4l_standard{index: 0}
	var cam bool
	switch err := ioctl_enumstd(d.fd, &s); err {
	case nil:
		cam = false
	case syscall.ENOTTY, syscall.EINVAL:
		cam = true
	default:
		return DeviceInfo{}, err
	}

	info := DeviceInfo{
		Path:       d.path,
		DeviceName: c.card,
		BusInfo:    c.busInfo,
		DriverName: c.driver,
		DriverVersion: [3]int{
			int(c.version>>16) & 0xff,
			int(c.version>>8) & 0xff,
			int(c.version) & 0xff,
		},
		Camera: cam,
	}
	return info, nil
}

// TurnOn initiates a capture session with the device. It may fail with
// ErrUnsupported. While the device is turned on, its configuration cannot be
// changed.
func (d *Device) TurnOn() error {
	// Switch to progressive format and reset the colorspace to device default.
	f := v4l_format_pix{typ: v4l_bufTypeVideoCapture}
	if err := ioctl_gFmt_pix(d.fd, &f); err != nil {
		return err
	}
	f.fmt.field = v4l_fieldNone
	f.fmt.colorspace = v4l_colorspaceDefault
	f.fmt.priv = 0
	if err := ioctl_sFmt_pix(d.fd, &f); err != nil {
		return err
	}

	// Reset cropping.
	cc := v4l_cropcap{typ: v4l_bufTypeVideoCapture}
	switch err := ioctl_cropcap(d.fd, &cc); err {
	case nil:
		c := v4l_crop{
			typ: v4l_bufTypeVideoCapture,
			c:   cc.defrect,
		}
		switch err := ioctl_sCrop(d.fd, &c); err {
		case nil:
			// Success.
		case syscall.ENOTTY, syscall.EINVAL:
			// VIDIOC_S_CROP unsupported.
		default:
			return err
		}
	case syscall.ENOTTY:
		// No support for cropping. That's okay.
	default:
		return err
	}

	// Allocate buffers.
	if err := d.allocBuffers(4); err != nil {
		return err
	}

	// Start streaming I/O.
	if err := ioctl_streamon(d.fd, v4l_bufTypeVideoCapture); err != nil {
		d.freeBuffers()
		return err
	}
	return nil
}

// TurnOff ends the capture session in progress. It does not close the device,
// so it can be reused for another session.
func (d *Device) TurnOff() {
	ioctl_streamoff(d.fd, v4l_bufTypeVideoCapture)
	d.freeBuffers()
}

// allocBuffers allocates n buffers in device memory, and mmaps and queues them.
func (d *Device) allocBuffers(n int) error {
	// Request buffers.
	rb := v4l_requestbuffers{
		count:  uint32(n),
		typ:    v4l_bufTypeVideoCapture,
		memory: v4l_memoryMmap,
	}
	if err := ioctl_reqbufs(d.fd, &rb); err != nil {
		if err == syscall.EINVAL {
			// Memory-mapped I/O method unsupported.
			err = ErrUnsupported
		}
		return err
	}
	if rb.count == 0 {
		return Error("out of device memory")
	}

	// Map and enqueue the buffers.
	d.buffers = make([][]byte, 0, rb.count)
	for i := 0; i < cap(d.buffers); i++ {
		b := v4l_buffer{
			index:  uint32(i),
			typ:    v4l_bufTypeVideoCapture,
			memory: v4l_memoryMmap,
		}
		if err := ioctl_querybuf(d.fd, &b); err != nil {
			d.freeBuffers()
			return err
		}
		buf, err := syscall.Mmap(d.fd, int64(b.offset), int(b.length),
			syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			d.freeBuffers()
			return err
		}
		d.buffers = append(d.buffers, buf)
		if err := ioctl_qbuf(d.fd, &b); err != nil {
			d.freeBuffers()
			return err
		}
	}

	return nil
}

// freeBuffers munmaps and frees any buffers allocated in device memory, and
// removes all pointers to them.
func (d *Device) freeBuffers() {
	d.bufIndex = noBuffer
	for i := range d.buffers {
		syscall.Munmap(d.buffers[i])
		d.buffers[i] = nil
	}
	d.buffers = nil
	rb := v4l_requestbuffers{
		count:  0,
		typ:    v4l_bufTypeVideoCapture,
		memory: v4l_memoryMmap,
	}
	ioctl_reqbufs(d.fd, &rb)
}

// Capture grabs the next frame, and returns a new Buffer holding the raw image
// data. The device must be turned on for Capture to succeed. A call to Capture
// may render the contents of previously captured buffers unavailable.
func (d *Device) Capture() (*Buffer, error) {
	d.nCaptures++

	// Enqueue the old buffer (if any).
	if d.bufIndex != noBuffer {
		b := v4l_buffer{
			typ:    v4l_bufTypeVideoCapture,
			memory: v4l_memoryMmap,
			index:  d.bufIndex,
		}
		d.bufIndex = noBuffer
		if err := ioctl_qbuf(d.fd, &b); err != nil {
			return nil, err
		}
	}

	// Dequeue a new buffer.
	b := v4l_buffer{
		typ:    v4l_bufTypeVideoCapture,
		memory: v4l_memoryMmap,
	}
	if err := ioctl_dqbuf(d.fd, &b); err != nil {
		return nil, err
	}
	d.buffers[b.index] = d.buffers[b.index][:b.bytesused]
	d.bufIndex = b.index

	return &Buffer{d.device, d.nCaptures, 0, b.sequence}, nil
}

// GetConfig returns the current configuration of the device.
func (d *Device) GetConfig() (DeviceConfig, error) {
	// Get format.
	f := v4l_format_pix{typ: v4l_bufTypeVideoCapture}
	if err := ioctl_gFmt_pix(d.fd, &f); err != nil {
		return DeviceConfig{}, err
	}

	// Get streaming parameters.
	p := v4l_streamparm_capture{typ: v4l_bufTypeVideoCapture}
	if err := ioctl_gParm_capture(d.fd, &p); err != nil {
		return DeviceConfig{}, err
	}

	cfg := DeviceConfig{
		Format: f.fmt.pixelformat,
		Width:  int(f.fmt.width),
		Height: int(f.fmt.height),
		FPS: Frac{
			p.parm.timeperframe.denominator,
			p.parm.timeperframe.numerator,
		},
	}
	return cfg, nil
}

// SetConfig configures the device according to cfg. The configuration actually
// applied may be different from what was requested, as drivers are allowed to
// adjust the parameters against hardware capabilities (or even completely
// ignore them). The configuration cannot be changed while the device is turned
// on.
func (d *Device) SetConfig(cfg DeviceConfig) error {
	// Set format.
	f := v4l_format_pix{
		typ: v4l_bufTypeVideoCapture,
		fmt: v4l_pixFormat{
			width:       uint32(cfg.Width),
			height:      uint32(cfg.Height),
			pixelformat: cfg.Format,
			field:       v4l_fieldNone,
			colorspace:  v4l_colorspaceDefault,
			priv:        0,
		},
	}
	if err := ioctl_sFmt_pix(d.fd, &f); err != nil {
		return err
	}

	// Set streaming parameters.
	cfg.FPS = cfg.FPS.Reduce()
	p := v4l_streamparm_capture{
		typ: v4l_bufTypeVideoCapture,
		parm: v4l_captureparm{
			timeperframe: v4l_fract{cfg.FPS.D, cfg.FPS.N},
		},
	}
	if err := ioctl_sParm_capture(d.fd, &p); err != nil {
		return err
	}

	return nil
}

// BufferInfo returns information about how image data is laid out in a buffer.
// For the same device configuration it always returns the same value.
func (d *Device) BufferInfo() (BufferInfo, error) {
	f := v4l_format_pix{typ: v4l_bufTypeVideoCapture}
	if err := ioctl_gFmt_pix(d.fd, &f); err != nil {
		return BufferInfo{}, err
	}
	info := BufferInfo{
		BufferSize:  int(f.fmt.sizeimage),
		ImageStride: int(f.fmt.bytesperline),
	}
	return info, nil
}

// ListConfigs returns the configurations supported by the device.
func (d *Device) ListConfigs() ([]DeviceConfig, error) {
	var cfgs []DeviceConfig
	for fmt := 0; ; fmt++ {
		fd := v4l_fmtdesc{
			index: uint32(fmt),
			typ:   v4l_bufTypeVideoCapture,
		}
		if err := ioctl_enumFmt(d.fd, &fd); err != nil {
			if err != syscall.EINVAL {
				return nil, err
			}
			break
		}
		sizes, err := d.enumFrameSizes(fd.pixelformat)
		if err != nil {
			return nil, err
		}
		for _, sz := range sizes {
			ivals, err := d.enumFrameIvals(fd.pixelformat, sz.width, sz.height)
			if err != nil {
				return nil, err
			}
			for _, ival := range ivals {
				cfg := DeviceConfig{
					Format: fd.pixelformat,
					Width:  int(sz.width),
					Height: int(sz.height),
					FPS: Frac{
						ival.denominator,
						ival.numerator,
					}.Reduce(),
				}
				cfgs = append(cfgs, cfg)
			}
		}
	}

	// Some devices/drivers seem to return duplicates. Fix that.
	n := 0
	for i := 0; i < len(cfgs); i++ {
		dupe := false
		for j := 0; j < n && !dupe; j++ {
			dupe = cfgs[i] == cfgs[j]
		}
		if !dupe {
			cfgs[n] = cfgs[i]
			n++
		}
	}
	cfgs = cfgs[:n]

	return cfgs, nil
}

// enumFrameSizes returns the supported frame sizes. If the device does not
// enumerate a discrete set of frame sizes, then a few common ones within the
// supported range are returned.
func (d *Device) enumFrameSizes(fmt uint32) ([]v4l_frmsizeDiscrete, error) {
	var (
		sizes []v4l_frmsizeDiscrete
		fs    v4l_frmsizeenum
	)

	// Try to enumerate discrete frame sizes.
loop:
	for index := 0; ; index++ {
		fs = v4l_frmsizeenum{
			index:       uint32(index),
			pixelFormat: fmt,
		}
		if err := ioctl_enumFramesizes(d.fd, &fs); err != nil {
			if err != syscall.EINVAL {
				return nil, err
			}
			return sizes, nil
		}
		switch fs.typ {
		case v4l_frmsizeTypeDiscrete:
			sizes = append(sizes, fs.discrete)
		case v4l_frmsizeTypeContinuous, v4l_frmsizeTypeStepwise:
			break loop
		default:
			return nil, nil
		}
	}

	// Fall back to a default list.
	fss := fs.stepwise
	for _, fsd := range defaultFrameSizes {
		if fsd.width < fss.minWidth || fsd.width > fss.maxWidth ||
			(fsd.width-fss.minWidth)%fss.stepWidth != 0 {
			continue
		}
		if fsd.height < fss.minHeight || fsd.height > fss.maxHeight ||
			(fsd.height-fss.minHeight)%fss.stepHeight != 0 {
			continue
		}
		sizes = append(sizes, fsd)
	}

	sizes = append(sizes, v4l_frmsizeDiscrete{fss.maxWidth, fss.maxHeight})

	return sizes, nil
}

// defaultFrameSizes lists a few common resolutions.
var defaultFrameSizes = []v4l_frmsizeDiscrete{
	{160, 120},
	{176, 144},
	{320, 180},
	{320, 240},
	{352, 288},
	{640, 360},
	{640, 480},
	{800, 600},
	{960, 540},
	{1024, 768},
	{1280, 720},
	{1280, 960},
	{1600, 1200},
	{1920, 1080},
	{3840, 2160},
	{7680, 4320},
}

// enumFrameIvals returns the supported frame intervals. If the device does not
// enumerate a discrete set of frame intervals, then a few common ones within
// the supported range are returned.
func (d *Device) enumFrameIvals(fmt, w, h uint32) ([]v4l_fract, error) {
	var (
		ivals []v4l_fract
		fi    v4l_frmivalenum
	)

	// Try to enumerate discrete frame intervals.
loop:
	for index := 0; ; index++ {
		fi = v4l_frmivalenum{
			index:       uint32(index),
			pixelFormat: fmt,
			width:       w,
			height:      h,
		}
		if err := ioctl_enumFrameintervals(d.fd, &fi); err != nil {
			if err != syscall.EINVAL {
				return nil, err
			}
			return ivals, nil
		}
		switch fi.typ {
		case v4l_frmivalTypeDiscrete:
			ivals = append(ivals, fi.discrete)
		case v4l_frmivalTypeContinuous, v4l_frmivalTypeStepwise:
			break loop
		default:
			return nil, nil
		}
	}

	// Fall back to a default list.
	var (
		fis   = fi.stepwise
		min_n = fis.min.numerator
		min_d = fis.min.denominator
		max_n = fis.max.numerator
		max_d = fis.max.denominator
	)
	for _, fid := range defaultFrameIntervals {
		if fid.numerator*min_d < min_n*fid.denominator ||
			fid.numerator*max_d > max_n*fid.denominator {
			continue
		}
		ivals = append(ivals, fid)
	}

	ivals = append(ivals, fis.min)

	return ivals, nil
}

// defaultFrameIntervals lists a few common frame intervals.
var defaultFrameIntervals = []v4l_fract{
	{1, 5},
	{1, 10},
	{1, 15},
	{1, 25},
	{1, 30},
	{1, 60},
}

// ControlInfo returns information about a control.
func (d *Device) ControlInfo(cid uint32) (ControlInfo, error) {
	info, err := d.controlInfo(cid)
	if err == errBadControl {
		// Pretend the control does not exist.
		err = syscall.EINVAL
	}
	return info, err
}

// ListControls returns the ControlInfo for every control the device has.
func (d *Device) ListControls() ([]ControlInfo, error) {
	var (
		lastCID uint32
		infos   []ControlInfo
	)

	for {
		info, err := d.controlInfo(lastCID | v4l_ctrlFlagNextCtrl)
		switch err {
		case nil:
			infos = append(infos, info)
			lastCID = info.CID
		case errBadControl:
			// Pretend the control does not exist.
			lastCID = info.CID
		case syscall.EINVAL:
			if lastCID == 0 {
				// No support for v4l_ctrlFlagNextCtrl.
				// Fall back to legacy method.
				return d.listControlsLegacy()
			}
			return infos, nil
		default:
			return nil, err
		}
	}
}

// listControlsLegacy enumerates all controls the device has by querying them
// one-by-one rather than using the v4l_ctrlFlagNextCtrl flag.
func (d *Device) listControlsLegacy() ([]ControlInfo, error) {
	var infos []ControlInfo

	// Standard controls.
	for cid := uint32(v4l_cidBase); cid < v4l_cidLastp1; cid++ {
		info, err := d.controlInfo(cid)
		if err != nil {
			if err == syscall.EINVAL || err == errBadControl {
				continue
			}
			return nil, err
		}
		infos = append(infos, info)
	}

	// Custom controls.
	for cid := uint32(v4l_cidPrivateBase); ; cid++ {
		info, err := d.controlInfo(cid)
		if err != nil {
			if err == syscall.EINVAL {
				break
			}
			if err == errBadControl {
				continue
			}
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// controlInfo returns information about a control. For disabled contorls and
// contorls of unsupported type it fails with errBadControl.
func (d *device) controlInfo(cid uint32) (ControlInfo, error) {
	qc := v4l_queryctrl{id: cid}
	if err := ioctl_queryctrl(d.fd, &qc); err != nil {
		return ControlInfo{}, err
	}

	info := ControlInfo{
		CID:     qc.id,
		Name:    qc.name,
		Min:     qc.minimum,
		Max:     qc.maximum,
		Step:    qc.step,
		Default: qc.defaultValue,
	}

	if qc.flags&v4l_ctrlFlagDisabled != 0 {
		return info, errBadControl
	}

	switch qc.typ {
	case v4l_ctrlTypeInteger:
		info.Type = "int"
	case v4l_ctrlTypeBoolean:
		info.Type = "bool"
	case v4l_ctrlTypeMenu:
		info.Type = "enum"
	case v4l_ctrlTypeButton:
		info.Type = "button"
	default:
		return info, errBadControl
	}

	if qc.typ == v4l_ctrlTypeMenu {
		for i := qc.minimum; i <= qc.maximum; i++ {
			qm := v4l_querymenu{
				id:    qc.id,
				index: uint32(i),
			}
			if err := ioctl_querymenu(d.fd, &qm); err != nil {
				if err == syscall.EINVAL {
					continue
				}
				return ControlInfo{}, err
			}
			opt := struct {
				Value int32
				Name  string
			}{
				Value: int32(qm.index),
				Name:  qm.name,
			}
			info.Options = append(info.Options, opt)
		}
	}

	return info, nil
}

// GetControl returns the current value of a control.
func (d *Device) GetControl(cid uint32) (int32, error) {
	c := v4l_control{id: cid}
	if err := ioctl_gCtrl(d.fd, &c); err != nil {
		return 0, err
	}
	return c.value, nil
}

// SetControl sets the value of a control.
func (d *Device) SetControl(cid uint32, value int32) error {
	c := v4l_control{
		id:    cid,
		value: value,
	}
	if err := ioctl_sCtrl(d.fd, &c); err != nil {
		return err
	}
	return nil
}
