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

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

var src = &Image{
	Pix: []byte{
		//       +---+---+---+---+---+---+---+---+
		//  Y:   | 0 |     y     | 0 |     x     |
		//       +---+---+---+---+---+---+---+---+
		//
		//       +---+---+---+---+---+---+---+---+
		//  Cb:  | 0 |     y     | 1 |  x/2  | 0 |
		//       +---+---+---+---+---+---+---+---+
		//
		//       +---+---+---+---+---+---+---+---+
		//  Cr:  | 0 |     y     | 1 |  x/2  | 1 |
		//       +---+---+---+---+---+---+---+---+

		0x00, 0x08, 0x01, 0x09, 0x02, 0x0a, 0x03, 0x0b, 0x04, 0x0c, 0x05, 0x0d, 0x06, 0x0e, 0x07, 0x0f, 0xff, 0xff, 0xff,
		0x10, 0x18, 0x11, 0x19, 0x12, 0x1a, 0x13, 0x1b, 0x14, 0x1c, 0x15, 0x1d, 0x16, 0x1e, 0x17, 0x1f, 0xff, 0xff, 0xff,
		0x20, 0x28, 0x21, 0x29, 0x22, 0x2a, 0x23, 0x2b, 0x24, 0x2c, 0x25, 0x2d, 0x26, 0x2e, 0x27, 0x2f, 0xff, 0xff, 0xff,
		0x30, 0x38, 0x31, 0x39, 0x32, 0x3a, 0x33, 0x3b, 0x34, 0x3c, 0x35, 0x3d, 0x36, 0x3e, 0x37, 0x3f, 0xff, 0xff, 0xff,
		0x40, 0x48, 0x41, 0x49, 0x42, 0x4a, 0x43, 0x4b, 0x44, 0x4c, 0x45, 0x4d, 0x46, 0x4e, 0x47, 0x4f, 0xff, 0xff, 0xff,
		0x50, 0x58, 0x51, 0x59, 0x52, 0x5a, 0x53, 0x5b, 0x54, 0x5c, 0x55, 0x5d, 0x56, 0x5e, 0x57, 0x5f, 0xff, 0xff, 0xff,
		0x60, 0x68, 0x61, 0x69, 0x62, 0x6a, 0x63, 0x6b, 0x64, 0x6c, 0x65, 0x6d, 0x66, 0x6e, 0x67, 0x6f, 0xff, 0xff, 0xff,
		0x70, 0x78, 0x71, 0x79, 0x72, 0x7a, 0x73, 0x7b, 0x74, 0x7c, 0x75, 0x7d, 0x76, 0x7e, 0x77, 0x7f, 0xff, 0xff, 0xff,
	},
	Stride: 19,
	Rect: image.Rectangle{
		Min: image.Point{10, 20},
		Max: image.Point{18, 28},
	},
}

func TestNew(t *testing.T) {
	// This test never fails explicitly. If it doesn't panic, it's okay.
	for L := 0; L < 4; L++ {
		for T := 0; T < 4; T++ {
			for R := 8; R < 12; R++ {
				for B := 8; B < 12; B++ {
					img := New(image.Rect(L, T, R, B))
					for y := T; y < B; y++ {
						for x := L; x < R; x++ {
							img.At(x, y)
						}
					}
				}
			}
		}
	}
}

func TestImage(t *testing.T) {
	for L := src.Rect.Min.X - 4; L < src.Rect.Max.X+4; L++ {
		for T := src.Rect.Min.Y - 4; T < src.Rect.Max.Y+4; T++ {
			for R := src.Rect.Min.X - 4; R < src.Rect.Max.X+4; R++ {
				for B := src.Rect.Min.Y - 4; B < src.Rect.Max.Y+4; B++ {
					r := image.Rect(L, T, R, B)
					s := src.SubImage(r)
					for y := src.Rect.Min.Y - 4; y < src.Rect.Max.Y+4; y++ {
						for x := src.Rect.Min.X - 4; x < src.Rect.Max.X+4; x++ {
							p := image.Pt(x, y)
							c0 := s.At(x, y)
							c1 := color.YCbCr{}
							if p.In(r) && p.In(src.Rect) {
								c1 = expectedAt(x, y)
							}
							if c0 != c1 {
								t.Errorf("got: %v, expected: %v (x=%d, y=%d, r=%v)\n",
									c0, c1, x, y, r)
								return
							}
						}
					}
				}
			}
		}
	}
}

func TestToRGBA(t *testing.T) {
	r0 := image.Rect(-2, -6, 14, 10)
	for L := r0.Min.X; L < r0.Max.X; L++ {
		for T := r0.Min.Y; T < r0.Max.Y; T++ {
			for R := r0.Min.X; R < r0.Max.X; R++ {
				for B := r0.Min.Y; B < r0.Max.Y; B++ {
					for X := src.Rect.Min.X - 4; X < src.Rect.Max.X+4; X++ {
						for Y := src.Rect.Min.Y - 4; Y < src.Rect.Max.Y+4; Y++ {
							r := image.Rect(L, T, R, B)
							p := image.Pt(X, Y)
							dst0 := image.NewRGBA(r0)
							dst1 := image.NewRGBA(r0)
							ToRGBA(dst0, r, src, p)
							draw.Draw(dst1, r, src, p, draw.Over)
							for x := r0.Min.X; x < r0.Max.X; x++ {
								for y := r0.Min.Y; y < r0.Max.Y; y++ {
									c0 := dst0.At(x, y)
									c1 := dst1.At(x, y)
									if c0 != c1 {
										t.Errorf("got: %v, expected: %v (x=%d, y=%d, r=%v, p=%v)\n",
											c0, c1, x, y, r, p)
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestToGray(t *testing.T) {
	r0 := image.Rect(-2, -6, 14, 10)
	for L := r0.Min.X; L < r0.Max.X; L++ {
		for T := r0.Min.Y; T < r0.Max.Y; T++ {
			for R := r0.Min.X; R < r0.Max.X; R++ {
				for B := r0.Min.Y; B < r0.Max.Y; B++ {
					for X := src.Rect.Min.X - 4; X < src.Rect.Max.X+4; X++ {
						for Y := src.Rect.Min.Y - 4; Y < src.Rect.Max.Y+4; Y++ {
							r := image.Rect(L, T, R, B)
							p := image.Pt(X, Y)
							dst0 := image.NewGray(r0)
							dst1 := image.NewGray(r0)
							ToGray(dst0, r, src, p)
							draw.Draw(dst1, r, src, p, draw.Over)
							for x := r0.Min.X; x < r0.Max.X; x++ {
								for y := r0.Min.Y; y < r0.Max.Y; y++ {
									c0 := dst0.At(x, y)
									c1 := dst1.At(x, y)
									if c0 != c1 {
										t.Errorf("got: %v, expected: %v (x=%d, y=%d, r=%v, p=%v)\n",
											c0, c1, x, y, r, p)
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestToYCbCr(t *testing.T) {
	r0 := image.Rect(-2, -6, 14, 10)
	for L := r0.Min.X; L < r0.Max.X; L++ {
		for T := r0.Min.Y; T < r0.Max.Y; T++ {
			for R := r0.Min.X; R < r0.Max.X; R++ {
				for B := r0.Min.Y; B < r0.Max.Y; B++ {
					for X := src.Rect.Min.X - 4; X < src.Rect.Max.X+4; X++ {
						for Y := src.Rect.Min.Y - 4; Y < src.Rect.Max.Y+4; Y++ {
							r := image.Rect(L, T, R, B)
							p := image.Pt(X, Y)
							d := p.Sub(r.Min)
							dst := image.NewYCbCr(r0, image.YCbCrSubsampleRatio444)
							ToYCbCr(dst, r, src, p)
							for x := r0.Min.X; x < r0.Max.X; x++ {
								for y := r0.Min.Y; y < r0.Max.Y; y++ {
									c0 := dst.At(x, y)
									c1 := color.YCbCr{}
									if image.Pt(x, y).In(r) && image.Pt(x, y).Add(d).In(src.Rect) {
										c1 = expectedAt(x+d.X, y+d.Y)
									}
									if c0 != c1 {
										t.Errorf("got: %v, expected: %v (x=%d, y=%d, r=%v, p=%v)\n",
											c0, c1, x, y, r, p)
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func expectedAt(x, y int) color.YCbCr {
	x -= src.Rect.Min.X
	y -= src.Rect.Min.Y
	return color.YCbCr{
		Y:  byte(y<<4 + x),
		Cb: byte(y<<4 + (x&^1 + 8)),
		Cr: byte(y<<4 + (x&^1 + 9)),
	}
}
