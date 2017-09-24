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

package yuyv

import "image"

// ToRGBA aligns r.Min in dst with p in src, and draws the part of src visible
// through r over src.
//
// It's several times faster than the image/draw package.
func ToRGBA(dst *image.RGBA, r image.Rectangle, src *Image, p image.Point) {
	v := p.Sub(r.Min)
	r = r.Intersect(dst.Rect).Intersect(src.Rect.Sub(v))
	p = r.Min.Add(v)
	if r.Empty() {
		return
	}
	for y := 0; y < r.Dy(); y++ {
		s := src.Pix[src.PixPairOffset(p.X, p.Y+y):]
		d := dst.Pix[dst.PixOffset(r.Min.X, r.Min.Y+y):]
		n := r.Dx()
		if p.X&1 != 0 {
			cb := int32(s[1]) - 128
			y2 := int32(s[2]) * 0x010101
			cr := int32(s[3]) - 128
			d[0] = clamp(y2 + 91881*cr)
			d[1] = clamp(y2 - 22554*cb - 46802*cr)
			d[2] = clamp(y2 + 116130*cb)
			d[3] = 255
			s = s[4:]
			d = d[4:]
			n--
		}
		for x := 0; x < n; x += 2 {
			y1 := int32(s[2*x+0]) * 0x010101
			cb := int32(s[2*x+1]) - 128
			y2 := int32(s[2*x+2]) * 0x010101
			cr := int32(s[2*x+3]) - 128
			d[4*x+0] = clamp(y1 + 91881*cr)
			d[4*x+1] = clamp(y1 - 22554*cb - 46802*cr)
			d[4*x+2] = clamp(y1 + 116130*cb)
			d[4*x+3] = 255
			if x < n-1 {
				d[4*x+4] = clamp(y2 + 91881*cr)
				d[4*x+5] = clamp(y2 - 22554*cb - 46802*cr)
				d[4*x+6] = clamp(y2 + 116130*cb)
				d[4*x+7] = 255
			}
		}
	}
}

func clamp(x int32) uint8 {
	if uint32(x)&0xff000000 == 0 {
		return uint8(x >> 16)
	} else {
		return uint8(^(x >> 31))
	}
}

// ToGray aligns r.Min in dst with p in src, and draws the part of src visible
// through r over src.
//
// It's several times faster than the image/draw package.
func ToGray(dst *image.Gray, r image.Rectangle, src *Image, p image.Point) {
	v := p.Sub(r.Min)
	r = r.Intersect(dst.Rect).Intersect(src.Rect.Sub(v))
	p = r.Min.Add(v)
	if r.Empty() {
		return
	}
	for y := 0; y < r.Dy(); y++ {
		s := src.Pix[src.PixPairOffset(p.X, p.Y+y):]
		d := dst.Pix[dst.PixOffset(r.Min.X, r.Min.Y+y):]
		n := r.Dx()
		if p.X&1 != 0 {
			cb := int32(s[1]) - 128
			y2 := int32(s[2]) * 0x010101
			cr := int32(s[3]) - 128
			r := clamp2(y2 + 91881*cr)
			g := clamp2(y2 - 22554*cb - 46802*cr)
			b := clamp2(y2 + 116130*cb)
			d[0] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			s = s[4:]
			d = d[1:]
			n--
		}
		for x := 0; x < n; x += 2 {
			y1 := int32(s[2*x+0]) * 0x010101
			cb := int32(s[2*x+1]) - 128
			y2 := int32(s[2*x+2]) * 0x010101
			cr := int32(s[2*x+3]) - 128
			r := clamp2(y1 + 91881*cr)
			g := clamp2(y1 - 22554*cb - 46802*cr)
			b := clamp2(y1 + 116130*cb)
			d[x] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			if x < n-1 {
				r := clamp2(y2 + 91881*cr)
				g := clamp2(y2 - 22554*cb - 46802*cr)
				b := clamp2(y2 + 116130*cb)
				d[x+1] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			}
		}
	}
}

func clamp2(x int32) int32 {
	if uint32(x)&0xff000000 == 0 {
		return (x >> 8) & 0xffff
	} else {
		return ^(x >> 31) & 0xffff
	}
}

// ToYCbCr aligns r.Min in dst with p in src, and draws the part of src visible
// through r over src.
//
// Panics if the subsample ratio of dst is not 4:4:4.
func ToYCbCr(dst *image.YCbCr, r image.Rectangle, src *Image, p image.Point) {
	if dst.SubsampleRatio != image.YCbCrSubsampleRatio444 {
		panic("subsample ratio must be 4:4:4")
	}
	v := p.Sub(r.Min)
	r = r.Intersect(dst.Rect).Intersect(src.Rect.Sub(v))
	p = r.Min.Add(v)
	if r.Empty() {
		return
	}
	for y := 0; y < r.Dy(); y++ {
		s := src.Pix[src.PixPairOffset(p.X, p.Y+y):]
		dy := dst.Y[dst.YOffset(r.Min.X, r.Min.Y+y):]
		dcb := dst.Cb[dst.COffset(r.Min.X, r.Min.Y+y):]
		dcr := dst.Cr[dst.COffset(r.Min.X, r.Min.Y+y):]
		n := r.Dx()
		if p.X&1 != 0 {
			dy[0], dcb[0], dcr[0] = s[2], s[1], s[3]
			s = s[4:]
			dy = dy[1:]
			dcb = dcb[1:]
			dcr = dcr[1:]
			n--
		}
		for x := 0; x < n; x += 2 {
			dy[x], dcb[x], dcr[x] = s[2*x], s[2*x+1], s[2*x+3]
			if x < n-1 {
				dy[x+1], dcb[x+1], dcr[x+1] = s[2*x+2], s[2*x+1], s[2*x+3]
			}
		}
	}
}
