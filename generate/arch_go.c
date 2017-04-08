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

#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#include <linux/videodev2.h>

int main() {
	printf("// +build linux\n");
	printf("// +build ");
	fflush(stdout);
	system("go env GOHOSTARCH");
	printf("\n");

	printf("/////////////////////////////////////////////////////\n");
	printf("//                                                 //\n");
	printf("//  !!! THIS IS A GENERATED FILE, DO NOT EDIT !!!  //\n");
	printf("//                                                 //\n");
	printf("/////////////////////////////////////////////////////\n");
	printf("\n");

	printf("package v4l\n\n");

	printf("const (\n");
	printf("\tvidioc_querycap           = 0x%08llx\n", (long long unsigned) VIDIOC_QUERYCAP);
	printf("\tvidioc_gFmt               = 0x%08llx\n", (long long unsigned) VIDIOC_G_FMT);
	printf("\tvidioc_sFmt               = 0x%08llx\n", (long long unsigned) VIDIOC_S_FMT);
	printf("\tvidioc_gParm              = 0x%08llx\n", (long long unsigned) VIDIOC_G_PARM);
	printf("\tvidioc_sParm              = 0x%08llx\n", (long long unsigned) VIDIOC_S_PARM);
	printf("\tvidioc_reqbufs            = 0x%08llx\n", (long long unsigned) VIDIOC_REQBUFS);
	printf("\tvidioc_querybuf           = 0x%08llx\n", (long long unsigned) VIDIOC_QUERYBUF);
	printf("\tvidioc_qbuf               = 0x%08llx\n", (long long unsigned) VIDIOC_QBUF);
	printf("\tvidioc_dqbuf              = 0x%08llx\n", (long long unsigned) VIDIOC_DQBUF);
	printf("\tvidioc_streamon           = 0x%08llx\n", (long long unsigned) VIDIOC_STREAMON);
	printf("\tvidioc_streamoff          = 0x%08llx\n", (long long unsigned) VIDIOC_STREAMOFF);
	printf("\tvidioc_cropcap            = 0x%08llx\n", (long long unsigned) VIDIOC_CROPCAP);
	printf("\tvidioc_sCrop              = 0x%08llx\n", (long long unsigned) VIDIOC_S_CROP);
	printf("\tvidioc_enumstd            = 0x%08llx\n", (long long unsigned) VIDIOC_ENUMSTD);
	printf("\tvidioc_enumFmt            = 0x%08llx\n", (long long unsigned) VIDIOC_ENUM_FMT);
	printf("\tvidioc_enumFramesizes     = 0x%08llx\n", (long long unsigned) VIDIOC_ENUM_FRAMESIZES);
	printf("\tvidioc_enumFrameintervals = 0x%08llx\n", (long long unsigned) VIDIOC_ENUM_FRAMEINTERVALS);
	printf("\tvidioc_queryctrl          = 0x%08llx\n", (long long unsigned) VIDIOC_QUERYCTRL);
	printf("\tvidioc_querymenu          = 0x%08llx\n", (long long unsigned) VIDIOC_QUERYMENU);
	printf("\tvidioc_gCtrl              = 0x%08llx\n", (long long unsigned) VIDIOC_G_CTRL);
	printf("\tvidioc_sCtrl              = 0x%08llx\n", (long long unsigned) VIDIOC_S_CTRL);
	printf(")\n\n");

	printf("const (\n");
	printf("\tsize_capability     = %llu\n", (long long unsigned) sizeof(struct v4l2_capability));
	printf("\tsize_format         = %llu\n", (long long unsigned) sizeof(struct v4l2_format));
	printf("\tsize_streamparm     = %llu\n", (long long unsigned) sizeof(struct v4l2_streamparm));
	printf("\tsize_requestbuffers = %llu\n", (long long unsigned) sizeof(struct v4l2_requestbuffers));
	printf("\tsize_buffer         = %llu\n", (long long unsigned) sizeof(struct v4l2_buffer));
	printf("\tsize_int            = %llu\n", (long long unsigned) sizeof(int));
	printf("\tsize_cropcap        = %llu\n", (long long unsigned) sizeof(struct v4l2_cropcap));
	printf("\tsize_crop           = %llu\n", (long long unsigned) sizeof(struct v4l2_crop));
	printf("\tsize_standard       = %llu\n", (long long unsigned) sizeof(struct v4l2_standard));
	printf("\tsize_fmtdesc        = %llu\n", (long long unsigned) sizeof(struct v4l2_fmtdesc));
	printf("\tsize_frmsizeenum    = %llu\n", (long long unsigned) sizeof(struct v4l2_frmsizeenum));
	printf("\tsize_frmivalenum    = %llu\n", (long long unsigned) sizeof(struct v4l2_frmivalenum));
	printf("\tsize_queryctrl      = %llu\n", (long long unsigned) sizeof(struct v4l2_queryctrl));
	printf("\tsize_querymenu      = %llu\n", (long long unsigned) sizeof(struct v4l2_querymenu));
	printf("\tsize_control        = %llu\n", (long long unsigned) sizeof(struct v4l2_control));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_capability_driver       = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, driver));
	printf("\tsize_capability_driver       = %llu\n", (long long unsigned) sizeof((struct v4l2_capability){0}.driver));
	printf("\toffs_capability_card         = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, card));
	printf("\tsize_capability_card         = %llu\n", (long long unsigned) sizeof((struct v4l2_capability){0}.card));
	printf("\toffs_capability_busInfo      = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, bus_info));
	printf("\tsize_capability_busInfo      = %llu\n", (long long unsigned) sizeof((struct v4l2_capability){0}.bus_info));
	printf("\toffs_capability_version      = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, version));
	printf("\toffs_capability_capabilities = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, capabilities));
	printf("\toffs_capability_deviceCaps   = %llu\n", (long long unsigned) offsetof(struct v4l2_capability, device_caps));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_format_typ = %llu\n", (long long unsigned) offsetof(struct v4l2_format, type));
	printf("\toffs_format_fmt = %llu\n", (long long unsigned) offsetof(struct v4l2_format, fmt));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_pixFormat_width        = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, width));
	printf("\toffs_pixFormat_height       = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, height));
	printf("\toffs_pixFormat_pixelformat  = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, pixelformat));
	printf("\toffs_pixFormat_field        = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, field));
	printf("\toffs_pixFormat_bytesperline = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, bytesperline));
	printf("\toffs_pixFormat_sizeimage    = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, sizeimage));
	printf("\toffs_pixFormat_colorspace   = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, colorspace));
	printf("\toffs_pixFormat_priv         = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, priv));
	printf("\toffs_pixFormat_flags        = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, flags));
	printf("\toffs_pixFormat_ycbcrEnc     = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, ycbcr_enc));
	printf("\toffs_pixFormat_quantization = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, quantization));
	printf("\toffs_pixFormat_xferFunc     = %llu\n", (long long unsigned) offsetof(struct v4l2_pix_format, xfer_func));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_streamparm_typ  = %llu\n", (long long unsigned) offsetof(struct v4l2_streamparm, type));
	printf("\toffs_streamparm_parm = %llu\n", (long long unsigned) offsetof(struct v4l2_streamparm, parm));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_captureparm_capability   = %llu\n", (long long unsigned) offsetof(struct v4l2_captureparm, capability));
	printf("\toffs_captureparm_capturemode  = %llu\n", (long long unsigned) offsetof(struct v4l2_captureparm, capturemode));
	printf("\toffs_captureparm_timeperframe = %llu\n", (long long unsigned) offsetof(struct v4l2_captureparm, timeperframe));
	printf("\toffs_captureparm_extendedmode = %llu\n", (long long unsigned) offsetof(struct v4l2_captureparm, extendedmode));
	printf("\toffs_captureparm_readbuffers  = %llu\n", (long long unsigned) offsetof(struct v4l2_captureparm, readbuffers));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_requestbuffers_count  = %llu\n", (long long unsigned) offsetof(struct v4l2_requestbuffers, count));
	printf("\toffs_requestbuffers_typ    = %llu\n", (long long unsigned) offsetof(struct v4l2_requestbuffers, type));
	printf("\toffs_requestbuffers_memory = %llu\n", (long long unsigned) offsetof(struct v4l2_requestbuffers, memory));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_buffer_index     = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, index));
	printf("\toffs_buffer_typ       = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, type));
	printf("\toffs_buffer_bytesused = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, bytesused));
	printf("\toffs_buffer_flags     = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, flags));
	printf("\toffs_buffer_field     = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, field));
	printf("\toffs_buffer_timecode  = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, timecode));
	printf("\toffs_buffer_sequence  = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, sequence));
	printf("\toffs_buffer_memory    = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, memory));
	printf("\toffs_buffer_offset    = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, m.offset));
	printf("\toffs_buffer_length    = %llu\n", (long long unsigned) offsetof(struct v4l2_buffer, length));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_cropcap_typ         = %llu\n", (long long unsigned) offsetof(struct v4l2_cropcap, type));
	printf("\toffs_cropcap_bounds      = %llu\n", (long long unsigned) offsetof(struct v4l2_cropcap, bounds));
	printf("\toffs_cropcap_defrect     = %llu\n", (long long unsigned) offsetof(struct v4l2_cropcap, defrect));
	printf("\toffs_cropcap_pixelaspect = %llu\n", (long long unsigned) offsetof(struct v4l2_cropcap, pixelaspect));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_crop_typ = %llu\n", (long long unsigned) offsetof(struct v4l2_crop, type));
	printf("\toffs_crop_c   = %llu\n", (long long unsigned) offsetof(struct v4l2_crop, c));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_fract_numerator   = %llu\n", (long long unsigned) offsetof(struct v4l2_fract, numerator));
	printf("\toffs_fract_denominator = %llu\n", (long long unsigned) offsetof(struct v4l2_fract, denominator));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_timecode_typ      = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, type));
	printf("\toffs_timecode_flags    = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, flags));
	printf("\toffs_timecode_frames   = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, frames));
	printf("\toffs_timecode_seconds  = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, seconds));
	printf("\toffs_timecode_minutes  = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, minutes));
	printf("\toffs_timecode_hours    = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, hours));
	printf("\toffs_timecode_userbits = %llu\n", (long long unsigned) offsetof(struct v4l2_timecode, userbits));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_rect_left   = %llu\n", (long long unsigned) offsetof(struct v4l2_rect, left));
	printf("\toffs_rect_top    = %llu\n", (long long unsigned) offsetof(struct v4l2_rect, top));
	printf("\toffs_rect_width  = %llu\n", (long long unsigned) offsetof(struct v4l2_rect, width));
	printf("\toffs_rect_height = %llu\n", (long long unsigned) offsetof(struct v4l2_rect, height));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_standard_index       = %llu\n", (long long unsigned) offsetof(struct v4l2_standard, index));
	printf("\toffs_standard_id          = %llu\n", (long long unsigned) offsetof(struct v4l2_standard, id));
	printf("\toffs_standard_name        = %llu\n", (long long unsigned) offsetof(struct v4l2_standard, name));
	printf("\tsize_standard_name        = %llu\n", (long long unsigned) sizeof((struct v4l2_standard){0}.name));
	printf("\toffs_standard_frameperiod = %llu\n", (long long unsigned) offsetof(struct v4l2_standard, frameperiod));
	printf("\toffs_standard_framelines  = %llu\n", (long long unsigned) offsetof(struct v4l2_standard, framelines));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_fmtdesc_index       = %llu\n", (long long unsigned) offsetof(struct v4l2_fmtdesc, index));
	printf("\toffs_fmtdesc_typ         = %llu\n", (long long unsigned) offsetof(struct v4l2_fmtdesc, type));
	printf("\toffs_fmtdesc_flags       = %llu\n", (long long unsigned) offsetof(struct v4l2_fmtdesc, flags));
	printf("\toffs_fmtdesc_description = %llu\n", (long long unsigned) offsetof(struct v4l2_fmtdesc, description));
	printf("\tsize_fmtdesc_description = %llu\n", (long long unsigned) sizeof((struct v4l2_fmtdesc){0}.description));
	printf("\toffs_fmtdesc_pixelformat = %llu\n", (long long unsigned) offsetof(struct v4l2_fmtdesc, pixelformat));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_frmsizeenum_index       = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsizeenum, index));
	printf("\toffs_frmsizeenum_pixelFormat = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsizeenum, pixel_format));
	printf("\toffs_frmsizeenum_typ         = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsizeenum, type));
	printf("\toffs_frmsizeenum_discrete    = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsizeenum, discrete));
	printf("\toffs_frmsizeenum_stepwise    = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsizeenum, stepwise));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_frmsizeDiscrete_width  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_discrete, width));
	printf("\toffs_frmsizeDiscrete_height = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_discrete, height));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_frmsizeStepwise_minWidth   = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, min_width));
	printf("\toffs_frmsizeStepwise_maxWidth   = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, max_width));
	printf("\toffs_frmsizeStepwise_stepWidth  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, step_width));
	printf("\toffs_frmsizeStepwise_minHeight  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, min_height));
	printf("\toffs_frmsizeStepwise_maxHeight  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, max_height));
	printf("\toffs_frmsizeStepwise_stepHeight = %llu\n", (long long unsigned) offsetof(struct v4l2_frmsize_stepwise, step_height));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_frmivalenum_index       = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, index));
	printf("\toffs_frmivalenum_pixelFormat = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, pixel_format));
	printf("\toffs_frmivalenum_width       = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, width));
	printf("\toffs_frmivalenum_height      = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, height));
	printf("\toffs_frmivalenum_typ         = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, type));
	printf("\toffs_frmivalenum_discrete    = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, discrete));
	printf("\toffs_frmivalenum_stepwise    = %llu\n", (long long unsigned) offsetof(struct v4l2_frmivalenum, stepwise));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_frmivalStepwise_min  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmival_stepwise, min));
	printf("\toffs_frmivalStepwise_max  = %llu\n", (long long unsigned) offsetof(struct v4l2_frmival_stepwise, max));
	printf("\toffs_frmivalStepwise_step = %llu\n", (long long unsigned) offsetof(struct v4l2_frmival_stepwise, step));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_queryctrl_id           = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, id));
	printf("\toffs_queryctrl_typ          = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, type));
	printf("\toffs_queryctrl_name         = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, name));
	printf("\tsize_queryctrl_name         = %llu\n", (long long unsigned) sizeof((struct v4l2_queryctrl){0}.name));
	printf("\toffs_queryctrl_minimum      = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, minimum));
	printf("\toffs_queryctrl_maximum      = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, maximum));
	printf("\toffs_queryctrl_step         = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, step));
	printf("\toffs_queryctrl_defaultValue = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, default_value));
	printf("\toffs_queryctrl_flags        = %llu\n", (long long unsigned) offsetof(struct v4l2_queryctrl, flags));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_querymenu_id    = %llu\n", (long long unsigned) offsetof(struct v4l2_querymenu, id));
	printf("\toffs_querymenu_index = %llu\n", (long long unsigned) offsetof(struct v4l2_querymenu, index));
	printf("\toffs_querymenu_name  = %llu\n", (long long unsigned) offsetof(struct v4l2_querymenu, name));
	printf("\tsize_querymenu_name  = %llu\n", (long long unsigned) sizeof((struct v4l2_querymenu){0}.name));
	printf("\toffs_querymenu_value = %llu\n", (long long unsigned) offsetof(struct v4l2_querymenu, value));
	printf(")\n\n");

	printf("const (\n");
	printf("\toffs_control_id    = %llu\n", (long long unsigned) offsetof(struct v4l2_control, id));
	printf("\toffs_control_value = %llu\n", (long long unsigned) offsetof(struct v4l2_control, value));
	printf(")\n");

	return 0;
}
