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

package v4l

import (
	"syscall"
	"unsafe"
)

// Constants.

const (
	v4l_capVideoCapture = 0x00000001
	v4l_capDeviceCaps   = 0x80000000
)

const (
	v4l_bufTypeVideoCapture = 1
)

const (
	v4l_colorspaceDefault = 0
)

const (
	v4l_fieldNone = 1
)

const (
	v4l_memoryMmap = 1
)

const (
	v4l_frmsizeTypeDiscrete   = 1
	v4l_frmsizeTypeContinuous = 2
	v4l_frmsizeTypeStepwise   = 3
)

const (
	v4l_frmivalTypeDiscrete   = 1
	v4l_frmivalTypeContinuous = 2
	v4l_frmivalTypeStepwise   = 3
)

const (
	v4l_ctrlFlagDisabled = 0x0001
	v4l_ctrlFlagNextCtrl = 0x80000000
)

const (
	v4l_ctrlTypeInteger     = 1
	v4l_ctrlTypeBoolean     = 2
	v4l_ctrlTypeMenu        = 3
	v4l_ctrlTypeButton      = 4
	v4l_ctrlTypeIntegerMenu = 9
)

const (
	v4l_cidBase        = 0x00980900
	v4l_cidLastp1      = 0x0098092b
	v4l_cidPrivateBase = 0x08000000
)

// Structs.

type v4l_capability struct {
	driver       string
	card         string
	busInfo      string
	version      uint32
	capabilities uint32
	deviceCaps   uint32
}

type v4l_format_pix struct {
	typ uint32
	fmt v4l_pixFormat
}

type v4l_pixFormat struct {
	width        uint32
	height       uint32
	pixelformat  uint32
	field        uint32
	bytesperline uint32
	sizeimage    uint32
	colorspace   uint32
	priv         uint32
	flags        uint32
	ycbcrEnc     uint32
	quantization uint32
	xferFunc     uint32
}

type v4l_streamparm_capture struct {
	typ  uint32
	parm v4l_captureparm
}

type v4l_captureparm struct {
	capability   uint32
	capturemode  uint32
	timeperframe v4l_fract
	extendedmode uint32
	readbuffers  uint32
}

type v4l_requestbuffers struct {
	count  uint32
	typ    uint32
	memory uint32
}

type v4l_buffer struct {
	index     uint32
	typ       uint32
	bytesused uint32
	flags     uint32
	field     uint32
	timecode  v4l_timecode
	sequence  uint32
	memory    uint32
	offset    uint32
	length    uint32
}

type v4l_int int32

type v4l_cropcap struct {
	typ         uint32
	bounds      v4l_rect
	defrect     v4l_rect
	pixelaspect v4l_fract
}

type v4l_crop struct {
	typ uint32
	c   v4l_rect
}

type v4l_fract struct {
	numerator   uint32
	denominator uint32
}

type v4l_timecode struct {
	typ      uint32
	flags    uint32
	frames   uint8
	seconds  uint8
	minutes  uint8
	hours    uint8
	userbits [4]uint8
}

type v4l_rect struct {
	left   int32
	top    int32
	width  uint32
	height uint32
}

type v4l_standard struct {
	index       uint32
	id          uint64
	name        string
	frameperiod v4l_fract
	framelines  uint32
}

type v4l_fmtdesc struct {
	index       uint32
	typ         uint32
	flags       uint32
	description string
	pixelformat uint32
}

type v4l_frmsizeenum struct {
	index       uint32
	pixelFormat uint32
	typ         uint32
	discrete    v4l_frmsizeDiscrete
	stepwise    v4l_frmsizeStepwise
}

type v4l_frmsizeDiscrete struct {
	width  uint32
	height uint32
}

type v4l_frmsizeStepwise struct {
	minWidth   uint32
	maxWidth   uint32
	stepWidth  uint32
	minHeight  uint32
	maxHeight  uint32
	stepHeight uint32
}

type v4l_frmivalenum struct {
	index       uint32
	pixelFormat uint32
	width       uint32
	height      uint32
	typ         uint32
	discrete    v4l_fract
	stepwise    v4l_frmivalStepwise
}

type v4l_frmivalStepwise struct {
	min  v4l_fract
	max  v4l_fract
	step v4l_fract
}

type v4l_queryctrl struct {
	id           uint32
	typ          uint32
	name         string
	minimum      int32
	maximum      int32
	step         int32
	defaultValue int32
	flags        uint32
}

type v4l_querymenu struct {
	id    uint32
	index uint32
	name  string
	value int64
}

type v4l_control struct {
	id    uint32
	value int32
}

// IOCTLs.

func ioctl_querycap(fd int, argp *v4l_capability) error {
	return ioctl(fd, vidioc_querycap, argp)
}

func ioctl_gFmt_pix(fd int, argp *v4l_format_pix) error {
	return ioctl(fd, vidioc_gFmt, argp)
}

func ioctl_sFmt_pix(fd int, argp *v4l_format_pix) error {
	return ioctl(fd, vidioc_sFmt, argp)
}

func ioctl_gParm_capture(fd int, argp *v4l_streamparm_capture) error {
	return ioctl(fd, vidioc_gParm, argp)
}

func ioctl_sParm_capture(fd int, argp *v4l_streamparm_capture) error {
	return ioctl(fd, vidioc_sParm, argp)
}

func ioctl_reqbufs(fd int, argp *v4l_requestbuffers) error {
	return ioctl(fd, vidioc_reqbufs, argp)
}

func ioctl_querybuf(fd int, argp *v4l_buffer) error {
	return ioctl(fd, vidioc_querybuf, argp)
}

func ioctl_qbuf(fd int, argp *v4l_buffer) error {
	return ioctl(fd, vidioc_qbuf, argp)
}

func ioctl_dqbuf(fd int, argp *v4l_buffer) error {
	return ioctl(fd, vidioc_dqbuf, argp)
}

func ioctl_streamon(fd int, typ v4l_int) error {
	return ioctl(fd, vidioc_streamon, &typ)
}

func ioctl_streamoff(fd int, typ v4l_int) error {
	return ioctl(fd, vidioc_streamoff, &typ)
}

func ioctl_cropcap(fd int, argp *v4l_cropcap) error {
	return ioctl(fd, vidioc_cropcap, argp)
}

func ioctl_sCrop(fd int, argp *v4l_crop) error {
	return ioctl(fd, vidioc_sCrop, argp)
}

func ioctl_enumstd(fd int, argp *v4l_standard) error {
	return ioctl(fd, vidioc_enumstd, argp)
}

func ioctl_enumFmt(fd int, argp *v4l_fmtdesc) error {
	return ioctl(fd, vidioc_enumFmt, argp)
}

func ioctl_enumFramesizes(fd int, argp *v4l_frmsizeenum) error {
	return ioctl(fd, vidioc_enumFramesizes, argp)
}

func ioctl_enumFrameintervals(fd int, argp *v4l_frmivalenum) error {
	return ioctl(fd, vidioc_enumFrameintervals, argp)
}

func ioctl_queryctrl(fd int, argp *v4l_queryctrl) error {
	return ioctl(fd, vidioc_queryctrl, argp)
}

func ioctl_querymenu(fd int, argp *v4l_querymenu) error {
	return ioctl(fd, vidioc_querymenu, argp)
}

func ioctl_gCtrl(fd int, argp *v4l_control) error {
	return ioctl(fd, vidioc_gCtrl, argp)
}

func ioctl_sCtrl(fd int, argp *v4l_control) error {
	return ioctl(fd, vidioc_sCtrl, argp)
}

func ioctl(fd int, request uint, argp ioctlArg) error {
	buf := make([]uint64, (argp.size()+7)/8)
	p := unsafe.Pointer(&buf[0])
	argp.put(p)
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(request), uintptr(p))
	if err != 0 {
		return err
	}
	argp.get(p)
	return nil
}

// Getters and putters.

type ioctlArg interface {
	get(unsafe.Pointer)
	put(unsafe.Pointer)
	size() int
}

func (p *v4l_capability) get(q unsafe.Pointer) {
	p.driver = getString(q, offs_capability_driver, size_capability_driver)
	p.card = getString(q, offs_capability_card, size_capability_card)
	p.busInfo = getString(q, offs_capability_busInfo, size_capability_busInfo)
	p.version = getUint32(q, offs_capability_version)
	p.capabilities = getUint32(q, offs_capability_capabilities)
	p.deviceCaps = getUint32(q, offs_capability_deviceCaps)
}

func (p *v4l_capability) put(q unsafe.Pointer) {
	putString(q, offs_capability_driver, size_capability_driver, p.driver)
	putString(q, offs_capability_card, size_capability_card, p.card)
	putString(q, offs_capability_busInfo, size_capability_busInfo, p.busInfo)
	putUint32(q, offs_capability_version, p.version)
	putUint32(q, offs_capability_capabilities, p.capabilities)
	putUint32(q, offs_capability_deviceCaps, p.deviceCaps)
}

func (p *v4l_capability) size() int {
	return size_capability
}

func (p *v4l_format_pix) get(q unsafe.Pointer) {
	p.typ = getUint32(q, offs_format_typ)
	p.fmt.get(unsafe.Pointer(uintptr(q) + offs_format_fmt))
}

func (p *v4l_format_pix) put(q unsafe.Pointer) {
	putUint32(q, offs_format_typ, p.typ)
	p.fmt.put(unsafe.Pointer(uintptr(q) + offs_format_fmt))
}

func (p *v4l_format_pix) size() int {
	return size_format
}

func (p *v4l_pixFormat) get(q unsafe.Pointer) {
	p.width = getUint32(q, offs_pixFormat_width)
	p.height = getUint32(q, offs_pixFormat_height)
	p.pixelformat = getUint32(q, offs_pixFormat_pixelformat)
	p.field = getUint32(q, offs_pixFormat_field)
	p.bytesperline = getUint32(q, offs_pixFormat_bytesperline)
	p.sizeimage = getUint32(q, offs_pixFormat_sizeimage)
	p.colorspace = getUint32(q, offs_pixFormat_colorspace)
	p.priv = getUint32(q, offs_pixFormat_priv)
	p.flags = getUint32(q, offs_pixFormat_flags)
	p.ycbcrEnc = getUint32(q, offs_pixFormat_ycbcrEnc)
	p.quantization = getUint32(q, offs_pixFormat_quantization)
	p.xferFunc = getUint32(q, offs_pixFormat_xferFunc)
}

func (p *v4l_pixFormat) put(q unsafe.Pointer) {
	putUint32(q, offs_pixFormat_width, p.width)
	putUint32(q, offs_pixFormat_height, p.height)
	putUint32(q, offs_pixFormat_pixelformat, p.pixelformat)
	putUint32(q, offs_pixFormat_field, p.field)
	putUint32(q, offs_pixFormat_bytesperline, p.bytesperline)
	putUint32(q, offs_pixFormat_sizeimage, p.sizeimage)
	putUint32(q, offs_pixFormat_colorspace, p.colorspace)
	putUint32(q, offs_pixFormat_priv, p.priv)
	putUint32(q, offs_pixFormat_flags, p.flags)
	putUint32(q, offs_pixFormat_ycbcrEnc, p.ycbcrEnc)
	putUint32(q, offs_pixFormat_quantization, p.quantization)
	putUint32(q, offs_pixFormat_xferFunc, p.xferFunc)
}

func (p *v4l_streamparm_capture) get(q unsafe.Pointer) {
	p.typ = getUint32(q, offs_streamparm_typ)
	p.parm.get(unsafe.Pointer(uintptr(q) + offs_streamparm_parm))
}

func (p *v4l_streamparm_capture) put(q unsafe.Pointer) {
	putUint32(q, offs_streamparm_typ, p.typ)
	p.parm.put(unsafe.Pointer(uintptr(q) + offs_streamparm_parm))
}

func (p *v4l_streamparm_capture) size() int {
	return size_streamparm
}

func (p *v4l_captureparm) get(q unsafe.Pointer) {
	p.capability = getUint32(q, offs_captureparm_capability)
	p.capturemode = getUint32(q, offs_captureparm_capturemode)
	p.timeperframe.get(unsafe.Pointer(uintptr(q) + offs_captureparm_timeperframe))
	p.extendedmode = getUint32(q, offs_captureparm_extendedmode)
	p.readbuffers = getUint32(q, offs_captureparm_readbuffers)
}

func (p *v4l_captureparm) put(q unsafe.Pointer) {
	putUint32(q, offs_captureparm_capability, p.capability)
	putUint32(q, offs_captureparm_capturemode, p.capturemode)
	p.timeperframe.put(unsafe.Pointer(uintptr(q) + offs_captureparm_timeperframe))
	putUint32(q, offs_captureparm_extendedmode, p.extendedmode)
	putUint32(q, offs_captureparm_readbuffers, p.readbuffers)
}

func (p *v4l_requestbuffers) get(q unsafe.Pointer) {
	p.count = getUint32(q, offs_requestbuffers_count)
	p.typ = getUint32(q, offs_requestbuffers_typ)
	p.memory = getUint32(q, offs_requestbuffers_memory)
}

func (p *v4l_requestbuffers) put(q unsafe.Pointer) {
	putUint32(q, offs_requestbuffers_count, p.count)
	putUint32(q, offs_requestbuffers_typ, p.typ)
	putUint32(q, offs_requestbuffers_memory, p.memory)
}

func (p *v4l_requestbuffers) size() int {
	return size_requestbuffers
}

func (p *v4l_buffer) get(q unsafe.Pointer) {
	p.index = getUint32(q, offs_buffer_index)
	p.typ = getUint32(q, offs_buffer_typ)
	p.bytesused = getUint32(q, offs_buffer_bytesused)
	p.flags = getUint32(q, offs_buffer_flags)
	p.field = getUint32(q, offs_buffer_field)
	p.timecode.get(unsafe.Pointer(uintptr(q) + offs_buffer_timecode))
	p.sequence = getUint32(q, offs_buffer_sequence)
	p.memory = getUint32(q, offs_buffer_memory)
	p.offset = getUint32(q, offs_buffer_offset)
	p.length = getUint32(q, offs_buffer_length)
}

func (p *v4l_buffer) put(q unsafe.Pointer) {
	putUint32(q, offs_buffer_index, p.index)
	putUint32(q, offs_buffer_typ, p.typ)
	putUint32(q, offs_buffer_bytesused, p.bytesused)
	putUint32(q, offs_buffer_flags, p.flags)
	putUint32(q, offs_buffer_field, p.field)
	p.timecode.put(unsafe.Pointer(uintptr(q) + offs_buffer_timecode))
	putUint32(q, offs_buffer_sequence, p.sequence)
	putUint32(q, offs_buffer_memory, p.memory)
	putUint32(q, offs_buffer_offset, p.offset)
	putUint32(q, offs_buffer_length, p.length)
}

func (p *v4l_buffer) size() int {
	return size_buffer
}

func (p *v4l_int) get(q unsafe.Pointer) {
	*p = v4l_int(getInt(q, 0))
}

func (p *v4l_int) put(q unsafe.Pointer) {
	putInt(q, 0, int32(*p))
}

func (p *v4l_int) size() int {
	return size_int
}

func (p *v4l_cropcap) get(q unsafe.Pointer) {
	p.typ = getUint32(q, offs_cropcap_typ)
	p.bounds.get(unsafe.Pointer(uintptr(q) + offs_cropcap_bounds))
	p.defrect.get(unsafe.Pointer(uintptr(q) + offs_cropcap_defrect))
	p.pixelaspect.get(unsafe.Pointer(uintptr(q) + offs_cropcap_pixelaspect))
}

func (p *v4l_cropcap) put(q unsafe.Pointer) {
	putUint32(q, offs_cropcap_typ, p.typ)
	p.bounds.put(unsafe.Pointer(uintptr(q) + offs_cropcap_bounds))
	p.defrect.put(unsafe.Pointer(uintptr(q) + offs_cropcap_defrect))
	p.pixelaspect.put(unsafe.Pointer(uintptr(q) + offs_cropcap_pixelaspect))
}

func (p *v4l_cropcap) size() int {
	return size_cropcap
}

func (p *v4l_crop) get(q unsafe.Pointer) {
	p.typ = getUint32(q, offs_crop_typ)
	p.c.get(unsafe.Pointer(uintptr(q) + offs_crop_c))
}

func (p *v4l_crop) put(q unsafe.Pointer) {
	putUint32(q, offs_crop_typ, p.typ)
	p.c.put(unsafe.Pointer(uintptr(q) + offs_crop_c))
}

func (p *v4l_crop) size() int {
	return size_crop
}

func (p *v4l_fract) get(q unsafe.Pointer) {
	p.numerator = getUint32(q, offs_fract_numerator)
	p.denominator = getUint32(q, offs_fract_denominator)
}

func (p *v4l_fract) put(q unsafe.Pointer) {
	putUint32(q, offs_fract_numerator, p.numerator)
	putUint32(q, offs_fract_denominator, p.denominator)
}

func (p *v4l_timecode) get(q unsafe.Pointer) {
	p.typ = getUint32(q, offs_timecode_typ)
	p.flags = getUint32(q, offs_timecode_flags)
	p.frames = getUint8(q, offs_timecode_frames)
	p.seconds = getUint8(q, offs_timecode_seconds)
	p.minutes = getUint8(q, offs_timecode_minutes)
	p.hours = getUint8(q, offs_timecode_hours)
	for i := range p.userbits {
		p.userbits[i] = getUint8(q, offs_timecode_userbits+i)
	}
}

func (p *v4l_timecode) put(q unsafe.Pointer) {
	putUint32(q, offs_timecode_typ, p.typ)
	putUint32(q, offs_timecode_flags, p.flags)
	putUint8(q, offs_timecode_frames, p.frames)
	putUint8(q, offs_timecode_seconds, p.seconds)
	putUint8(q, offs_timecode_minutes, p.minutes)
	putUint8(q, offs_timecode_hours, p.hours)
	for i := range p.userbits {
		putUint8(q, offs_timecode_userbits+i, p.userbits[i])
	}
}

func (p *v4l_rect) get(q unsafe.Pointer) {
	p.left = getInt32(q, offs_rect_left)
	p.top = getInt32(q, offs_rect_top)
	p.width = getUint32(q, offs_rect_width)
	p.height = getUint32(q, offs_rect_height)
}

func (p *v4l_rect) put(q unsafe.Pointer) {
	putInt32(q, offs_rect_left, p.left)
	putInt32(q, offs_rect_top, p.top)
	putUint32(q, offs_rect_width, p.width)
	putUint32(q, offs_rect_height, p.height)
}

func (p *v4l_standard) get(q unsafe.Pointer) {
	p.index = getUint32(q, offs_standard_index)
	p.id = getUint64(q, offs_standard_id)
	p.name = getString(q, offs_standard_name, size_standard_name)
	p.frameperiod.get(unsafe.Pointer(uintptr(q) + offs_standard_frameperiod))
	p.framelines = getUint32(q, offs_standard_framelines)
}

func (p *v4l_standard) put(q unsafe.Pointer) {
	putUint32(q, offs_standard_index, p.index)
	putUint64(q, offs_standard_id, p.id)
	putString(q, offs_standard_name, size_standard_name, p.name)
	p.frameperiod.put(unsafe.Pointer(uintptr(q) + offs_standard_frameperiod))
	putUint32(q, offs_standard_framelines, p.framelines)
}

func (p *v4l_standard) size() int {
	return size_standard
}

func (p *v4l_fmtdesc) get(q unsafe.Pointer) {
	p.index = getUint32(q, offs_fmtdesc_index)
	p.typ = getUint32(q, offs_fmtdesc_typ)
	p.flags = getUint32(q, offs_fmtdesc_flags)
	p.description = getString(q, offs_fmtdesc_description, size_fmtdesc_description)
	p.pixelformat = getUint32(q, offs_fmtdesc_pixelformat)
}

func (p *v4l_fmtdesc) put(q unsafe.Pointer) {
	putUint32(q, offs_fmtdesc_index, p.index)
	putUint32(q, offs_fmtdesc_typ, p.typ)
	putUint32(q, offs_fmtdesc_flags, p.flags)
	putString(q, offs_fmtdesc_description, size_fmtdesc_description, p.description)
	putUint32(q, offs_fmtdesc_pixelformat, p.pixelformat)

}

func (p *v4l_fmtdesc) size() int {
	return size_fmtdesc
}

func (p *v4l_frmsizeenum) get(q unsafe.Pointer) {
	p.index = getUint32(q, offs_frmsizeenum_index)
	p.pixelFormat = getUint32(q, offs_frmsizeenum_pixelFormat)
	p.typ = getUint32(q, offs_frmsizeenum_typ)
	p.discrete.get(unsafe.Pointer(uintptr(q) + offs_frmsizeenum_discrete))
	p.stepwise.get(unsafe.Pointer(uintptr(q) + offs_frmsizeenum_stepwise))
}

func (p *v4l_frmsizeenum) put(q unsafe.Pointer) {
	putUint32(q, offs_frmsizeenum_index, p.index)
	putUint32(q, offs_frmsizeenum_pixelFormat, p.pixelFormat)
	putUint32(q, offs_frmsizeenum_typ, p.typ)
	switch p.typ {
	case v4l_frmsizeTypeDiscrete:
		p.discrete.put(unsafe.Pointer(uintptr(q) + offs_frmsizeenum_discrete))
	case v4l_frmsizeTypeContinuous, v4l_frmsizeTypeStepwise:
		p.stepwise.put(unsafe.Pointer(uintptr(q) + offs_frmsizeenum_stepwise))
	}
}

func (p *v4l_frmsizeenum) size() int {
	return size_frmsizeenum
}

func (p *v4l_frmsizeDiscrete) get(q unsafe.Pointer) {
	p.width = getUint32(q, offs_frmsizeDiscrete_width)
	p.height = getUint32(q, offs_frmsizeDiscrete_height)
}

func (p *v4l_frmsizeDiscrete) put(q unsafe.Pointer) {
	putUint32(q, offs_frmsizeDiscrete_width, p.width)
	putUint32(q, offs_frmsizeDiscrete_height, p.height)
}

func (p *v4l_frmsizeStepwise) get(q unsafe.Pointer) {
	p.minWidth = getUint32(q, offs_frmsizeStepwise_minWidth)
	p.maxWidth = getUint32(q, offs_frmsizeStepwise_maxWidth)
	p.stepWidth = getUint32(q, offs_frmsizeStepwise_stepWidth)
	p.minHeight = getUint32(q, offs_frmsizeStepwise_minHeight)
	p.maxHeight = getUint32(q, offs_frmsizeStepwise_maxHeight)
	p.stepHeight = getUint32(q, offs_frmsizeStepwise_stepHeight)
}

func (p *v4l_frmsizeStepwise) put(q unsafe.Pointer) {
	putUint32(q, offs_frmsizeStepwise_minWidth, p.minWidth)
	putUint32(q, offs_frmsizeStepwise_maxWidth, p.maxWidth)
	putUint32(q, offs_frmsizeStepwise_stepWidth, p.stepWidth)
	putUint32(q, offs_frmsizeStepwise_minHeight, p.minHeight)
	putUint32(q, offs_frmsizeStepwise_maxHeight, p.maxHeight)
	putUint32(q, offs_frmsizeStepwise_stepHeight, p.stepHeight)
}

func (p *v4l_frmivalenum) get(q unsafe.Pointer) {
	p.index = getUint32(q, offs_frmivalenum_index)
	p.pixelFormat = getUint32(q, offs_frmivalenum_pixelFormat)
	p.width = getUint32(q, offs_frmivalenum_width)
	p.height = getUint32(q, offs_frmivalenum_height)
	p.typ = getUint32(q, offs_frmivalenum_typ)
	p.discrete.get(unsafe.Pointer(uintptr(q) + offs_frmivalenum_discrete))
	p.stepwise.get(unsafe.Pointer(uintptr(q) + offs_frmivalenum_stepwise))
}

func (p *v4l_frmivalenum) put(q unsafe.Pointer) {
	putUint32(q, offs_frmivalenum_index, p.index)
	putUint32(q, offs_frmivalenum_pixelFormat, p.pixelFormat)
	putUint32(q, offs_frmivalenum_width, p.width)
	putUint32(q, offs_frmivalenum_height, p.height)
	putUint32(q, offs_frmivalenum_typ, p.typ)
	switch p.typ {
	case v4l_frmivalTypeDiscrete:
		p.discrete.put(unsafe.Pointer(uintptr(q) + offs_frmivalenum_discrete))
	case v4l_frmivalTypeContinuous, v4l_frmivalTypeStepwise:
		p.stepwise.put(unsafe.Pointer(uintptr(q) + offs_frmivalenum_stepwise))
	}
}

func (p *v4l_frmivalenum) size() int {
	return size_frmivalenum
}

func (p *v4l_frmivalStepwise) get(q unsafe.Pointer) {
	p.min.get(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_min))
	p.max.get(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_max))
	p.step.get(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_step))
}

func (p *v4l_frmivalStepwise) put(q unsafe.Pointer) {
	p.min.put(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_min))
	p.max.put(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_max))
	p.step.put(unsafe.Pointer(uintptr(q) + offs_frmivalStepwise_step))
}

func (p *v4l_queryctrl) get(q unsafe.Pointer) {
	p.id = getUint32(q, offs_queryctrl_id)
	p.typ = getUint32(q, offs_queryctrl_typ)
	p.name = getString(q, offs_queryctrl_name, size_queryctrl_name)
	p.minimum = getInt32(q, offs_queryctrl_minimum)
	p.maximum = getInt32(q, offs_queryctrl_maximum)
	p.step = getInt32(q, offs_queryctrl_step)
	p.defaultValue = getInt32(q, offs_queryctrl_defaultValue)
	p.flags = getUint32(q, offs_queryctrl_flags)
}

func (p *v4l_queryctrl) put(q unsafe.Pointer) {
	putUint32(q, offs_queryctrl_id, p.id)
	putUint32(q, offs_queryctrl_typ, p.typ)
	putString(q, offs_queryctrl_name, size_queryctrl_name, p.name)
	putInt32(q, offs_queryctrl_minimum, p.minimum)
	putInt32(q, offs_queryctrl_maximum, p.maximum)
	putInt32(q, offs_queryctrl_step, p.step)
	putInt32(q, offs_queryctrl_defaultValue, p.defaultValue)
	putUint32(q, offs_queryctrl_flags, p.flags)
}

func (p *v4l_queryctrl) size() int {
	return size_queryctrl
}

func (p *v4l_querymenu) get(q unsafe.Pointer) {
	p.id = getUint32(q, offs_querymenu_id)
	p.index = getUint32(q, offs_querymenu_index)
	p.name = getString(q, offs_querymenu_name, size_querymenu_name)
	p.value = getInt64(q, offs_querymenu_value)
}

func (p *v4l_querymenu) put(q unsafe.Pointer) {
	putUint32(q, offs_querymenu_id, p.id)
	putUint32(q, offs_querymenu_index, p.index)
}

func (p *v4l_querymenu) size() int {
	return size_querymenu
}

func (p *v4l_control) get(q unsafe.Pointer) {
	p.id = getUint32(q, offs_control_id)
	p.value = getInt32(q, offs_control_value)
}

func (p *v4l_control) put(q unsafe.Pointer) {
	putUint32(q, offs_control_id, p.id)
	putInt32(q, offs_control_value, p.value)
}

func (p *v4l_control) size() int {
	return size_control
}

// Getters and putters for built-in types.

func getUint64(base unsafe.Pointer, offset int) uint64 {
	ptr := (*uint64)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	return *ptr
}

func putUint64(base unsafe.Pointer, offset int, value uint64) {
	ptr := (*uint64)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	*ptr = value
}

func getInt64(base unsafe.Pointer, offset int) int64 {
	ptr := (*int64)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	return *ptr
}

func putInt64(base unsafe.Pointer, offset int, value int64) {
	ptr := (*int64)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	*ptr = value
}

func getUint32(base unsafe.Pointer, offset int) uint32 {
	ptr := (*uint32)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	return *ptr
}

func putUint32(base unsafe.Pointer, offset int, value uint32) {
	ptr := (*uint32)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	*ptr = value
}

func getInt32(base unsafe.Pointer, offset int) int32 {
	ptr := (*int32)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	return *ptr
}

func putInt32(base unsafe.Pointer, offset int, value int32) {
	ptr := (*int32)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	*ptr = value
}

func getUint8(base unsafe.Pointer, offset int) uint8 {
	ptr := (*uint8)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	return *ptr
}

func putUint8(base unsafe.Pointer, offset int, value uint8) {
	ptr := (*uint8)(unsafe.Pointer(uintptr(base) + uintptr(offset)))
	*ptr = value
}

func getInt(base unsafe.Pointer, offset int) int64 {
	p := unsafe.Pointer(uintptr(base) + uintptr(offset))
	switch size_int {
	case 4:
		return int64(*(*uint32)(p))
	case 8:
		return *(*int64)(p)
	default:
		panic("bad int size")
	}
}

func putInt(base unsafe.Pointer, offset int, value int32) {
	p := unsafe.Pointer(uintptr(base) + uintptr(offset))
	switch size_int {
	case 4:
		*(*int32)(p) = value
	case 8:
		*(*int64)(p) = int64(value)
	default:
		panic("bad int size")
	}
}

func getString(base unsafe.Pointer, offset, maxLen int) string {
	buf := make([]byte, 0, maxLen)
	for i := 0; i < maxLen; i++ {
		ptr := (*byte)(unsafe.Pointer(uintptr(base) + uintptr(offset+i)))
		ch := *ptr
		if ch == 0 {
			break
		}
		buf = append(buf, ch)
	}
	return string(buf)
}

func putString(base unsafe.Pointer, offset, maxLen int, value string) {
	for i := 0; i < maxLen; i++ {
		ptr := (*byte)(unsafe.Pointer(uintptr(base) + uintptr(offset+i)))
		var ch byte
		if i < len(value) && i != maxLen-1 {
			ch = value[i]
		}
		*ptr = ch
	}
}
