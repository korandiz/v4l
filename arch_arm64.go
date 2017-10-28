// +build linux
// +build arm64

/////////////////////////////////////////////////////
//                                                 //
//  !!! THIS IS A GENERATED FILE, DO NOT EDIT !!!  //
//                                                 //
/////////////////////////////////////////////////////

package v4l

const (
	vidioc_querycap           = 0x80685600
	vidioc_gFmt               = 0xc0d05604
	vidioc_sFmt               = 0xc0d05605
	vidioc_gParm              = 0xc0cc5615
	vidioc_sParm              = 0xc0cc5616
	vidioc_reqbufs            = 0xc0145608
	vidioc_querybuf           = 0xc0585609
	vidioc_qbuf               = 0xc058560f
	vidioc_dqbuf              = 0xc0585611
	vidioc_streamon           = 0x40045612
	vidioc_streamoff          = 0x40045613
	vidioc_cropcap            = 0xc02c563a
	vidioc_sCrop              = 0x4014563c
	vidioc_enumstd            = 0xc0485619
	vidioc_enumFmt            = 0xc0405602
	vidioc_enumFramesizes     = 0xc02c564a
	vidioc_enumFrameintervals = 0xc034564b
	vidioc_queryctrl          = 0xc0445624
	vidioc_querymenu          = 0xc02c5625
	vidioc_gCtrl              = 0xc008561b
	vidioc_sCtrl              = 0xc008561c
)

const (
	size_capability     = 104
	size_format         = 208
	size_streamparm     = 204
	size_requestbuffers = 20
	size_buffer         = 88
	size_int            = 4
	size_cropcap        = 44
	size_crop           = 20
	size_standard       = 72
	size_fmtdesc        = 64
	size_frmsizeenum    = 44
	size_frmivalenum    = 52
	size_queryctrl      = 68
	size_querymenu      = 44
	size_control        = 8
)

const (
	offs_capability_driver       = 0
	size_capability_driver       = 16
	offs_capability_card         = 16
	size_capability_card         = 32
	offs_capability_busInfo      = 48
	size_capability_busInfo      = 32
	offs_capability_version      = 80
	offs_capability_capabilities = 84
	offs_capability_deviceCaps   = 88
)

const (
	offs_format_typ = 0
	offs_format_fmt = 8
)

const (
	offs_pixFormat_width        = 0
	offs_pixFormat_height       = 4
	offs_pixFormat_pixelformat  = 8
	offs_pixFormat_field        = 12
	offs_pixFormat_bytesperline = 16
	offs_pixFormat_sizeimage    = 20
	offs_pixFormat_colorspace   = 24
	offs_pixFormat_priv         = 28
	offs_pixFormat_flags        = 32
	offs_pixFormat_ycbcrEnc     = 36
	offs_pixFormat_quantization = 40
	offs_pixFormat_xferFunc     = 44
)

const (
	offs_streamparm_typ  = 0
	offs_streamparm_parm = 4
)

const (
	offs_captureparm_capability   = 0
	offs_captureparm_capturemode  = 4
	offs_captureparm_timeperframe = 8
	offs_captureparm_extendedmode = 16
	offs_captureparm_readbuffers  = 20
)

const (
	offs_requestbuffers_count  = 0
	offs_requestbuffers_typ    = 4
	offs_requestbuffers_memory = 8
)

const (
	offs_buffer_index     = 0
	offs_buffer_typ       = 4
	offs_buffer_bytesused = 8
	offs_buffer_flags     = 12
	offs_buffer_field     = 16
	offs_buffer_timecode  = 40
	offs_buffer_sequence  = 56
	offs_buffer_memory    = 60
	offs_buffer_offset    = 64
	offs_buffer_length    = 72
)

const (
	offs_cropcap_typ         = 0
	offs_cropcap_bounds      = 4
	offs_cropcap_defrect     = 20
	offs_cropcap_pixelaspect = 36
)

const (
	offs_crop_typ = 0
	offs_crop_c   = 4
)

const (
	offs_fract_numerator   = 0
	offs_fract_denominator = 4
)

const (
	offs_timecode_typ      = 0
	offs_timecode_flags    = 4
	offs_timecode_frames   = 8
	offs_timecode_seconds  = 9
	offs_timecode_minutes  = 10
	offs_timecode_hours    = 11
	offs_timecode_userbits = 12
)

const (
	offs_rect_left   = 0
	offs_rect_top    = 4
	offs_rect_width  = 8
	offs_rect_height = 12
)

const (
	offs_standard_index       = 0
	offs_standard_id          = 8
	offs_standard_name        = 16
	size_standard_name        = 24
	offs_standard_frameperiod = 40
	offs_standard_framelines  = 48
)

const (
	offs_fmtdesc_index       = 0
	offs_fmtdesc_typ         = 4
	offs_fmtdesc_flags       = 8
	offs_fmtdesc_description = 12
	size_fmtdesc_description = 32
	offs_fmtdesc_pixelformat = 44
)

const (
	offs_frmsizeenum_index       = 0
	offs_frmsizeenum_pixelFormat = 4
	offs_frmsizeenum_typ         = 8
	offs_frmsizeenum_discrete    = 12
	offs_frmsizeenum_stepwise    = 12
)

const (
	offs_frmsizeDiscrete_width  = 0
	offs_frmsizeDiscrete_height = 4
)

const (
	offs_frmsizeStepwise_minWidth   = 0
	offs_frmsizeStepwise_maxWidth   = 4
	offs_frmsizeStepwise_stepWidth  = 8
	offs_frmsizeStepwise_minHeight  = 12
	offs_frmsizeStepwise_maxHeight  = 16
	offs_frmsizeStepwise_stepHeight = 20
)

const (
	offs_frmivalenum_index       = 0
	offs_frmivalenum_pixelFormat = 4
	offs_frmivalenum_width       = 8
	offs_frmivalenum_height      = 12
	offs_frmivalenum_typ         = 16
	offs_frmivalenum_discrete    = 20
	offs_frmivalenum_stepwise    = 20
)

const (
	offs_frmivalStepwise_min  = 0
	offs_frmivalStepwise_max  = 8
	offs_frmivalStepwise_step = 16
)

const (
	offs_queryctrl_id           = 0
	offs_queryctrl_typ          = 4
	offs_queryctrl_name         = 8
	size_queryctrl_name         = 32
	offs_queryctrl_minimum      = 40
	offs_queryctrl_maximum      = 44
	offs_queryctrl_step         = 48
	offs_queryctrl_defaultValue = 52
	offs_queryctrl_flags        = 56
)

const (
	offs_querymenu_id    = 0
	offs_querymenu_index = 4
	offs_querymenu_name  = 8
	size_querymenu_name  = 32
	offs_querymenu_value = 8
)

const (
	offs_control_id    = 0
	offs_control_value = 4
)
