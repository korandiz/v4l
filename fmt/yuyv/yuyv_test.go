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
		//  Y:   |     x     |     y     | 0 | 0 |
		//       +---+---+---+---+---+---+---+---+
		//
		//       +---+---+---+---+---+---+---+---+
		//  Cb:  |  x/2  |     y     | 0 | 0 | 1 |
		//       +---+---+---+---+---+---+---+---+
		//
		//       +---+---+---+---+---+---+---+---+
		//  Cr:  |  x/2  |     y     | 1 | 1 | 1 |
		//       +---+---+---+---+---+---+---+---+
		0x00, 0x01, 0x20, 0x07, 0x40, 0x41, 0x60, 0x47, 0x80, 0x81, 0xa0, 0x87, 0xc0, 0xc1, 0xe0, 0xc7, 0xff, 0xff, 0xff,
		0x04, 0x09, 0x24, 0x0f, 0x44, 0x49, 0x64, 0x4f, 0x84, 0x89, 0xa4, 0x8f, 0xc4, 0xc9, 0xe4, 0xcf, 0xff, 0xff, 0xff,
		0x08, 0x11, 0x28, 0x17, 0x48, 0x51, 0x68, 0x57, 0x88, 0x91, 0xa8, 0x97, 0xc8, 0xd1, 0xe8, 0xd7, 0xff, 0xff, 0xff,
		0x0c, 0x19, 0x2c, 0x1f, 0x4c, 0x59, 0x6c, 0x5f, 0x8c, 0x99, 0xac, 0x9f, 0xcc, 0xd9, 0xec, 0xdf, 0xff, 0xff, 0xff,
		0x10, 0x21, 0x30, 0x27, 0x50, 0x61, 0x70, 0x67, 0x90, 0xa1, 0xb0, 0xa7, 0xd0, 0xe1, 0xf0, 0xe7, 0xff, 0xff, 0xff,
		0x14, 0x29, 0x34, 0x2f, 0x54, 0x69, 0x74, 0x6f, 0x94, 0xa9, 0xb4, 0xaf, 0xd4, 0xe9, 0xf4, 0xef, 0xff, 0xff, 0xff,
		0x18, 0x31, 0x38, 0x37, 0x58, 0x71, 0x78, 0x77, 0x98, 0xb1, 0xb8, 0xb7, 0xd8, 0xf1, 0xf8, 0xf7, 0xff, 0xff, 0xff,
		0x1c, 0x39, 0x3c, 0x3f, 0x5c, 0x79, 0x7c, 0x7f, 0x9c, 0xb9, 0xbc, 0xbf, 0xdc, 0xf9, 0xfc, 0xff, 0xff, 0xff, 0xff,
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	x = (x - src.Rect.Min.X) & 7
	y = (y - src.Rect.Min.Y) & 7
	return color.YCbCr{
		Y:  byte(x<<5 + y<<2),
		Cb: byte(x&^1<<5 + y<<3 + 1),
		Cr: byte(x&^1<<5 + y<<3 + 7),
	}
}

func TestRGBConversion(t *testing.T) {
	var (
		ycc  = New(image.Rect(0, 0, 2, 1))
		rgb1 = image.NewRGBA(ycc.Rect)
		rgb2 = image.NewRGBA(ycc.Rect)
	)
	for y := 0; y < 256; y++ {
		for cb := 0; cb < 256; cb++ {
			for cr := 0; cr < 256; cr++ {
				var (
					y  = uint8(y)
					cb = uint8(cb)
					cr = uint8(cr)
				)
				ycc.Pix[0] = y
				ycc.Pix[1] = cb
				ycc.Pix[2] = y
				ycc.Pix[3] = cr
				draw.Draw(rgb1, rgb1.Rect, ycc, ycc.Rect.Min, draw.Over)
				ToRGBA(rgb2, rgb2.Rect, ycc, ycc.Rect.Min)
				for i, v := range rgb2.Pix {
					if v != rgb1.Pix[i] {
						t.Errorf("got: %v, expected: %v (y=%d, cb=%d, cr=%d)\n",
							rgb2.Pix, rgb1.Pix, y, cb, cr)
						return
					}
				}
			}
		}
	}
}

func TestGrayConversion(t *testing.T) {
	var (
		ycc   = New(image.Rect(0, 0, 2, 1))
		gray1 = image.NewGray(ycc.Rect)
		gray2 = image.NewGray(ycc.Rect)
	)
	for y := 0; y < 256; y++ {
		for cb := 0; cb < 256; cb++ {
			for cr := 0; cr < 256; cr++ {
				var (
					y  = uint8(y)
					cb = uint8(cb)
					cr = uint8(cr)
				)
				ycc.Pix[0] = y
				ycc.Pix[1] = cb
				ycc.Pix[2] = y
				ycc.Pix[3] = cr
				draw.Draw(gray1, gray1.Rect, ycc, ycc.Rect.Min, draw.Over)
				ToGray(gray2, gray2.Rect, ycc, ycc.Rect.Min)
				for i, v := range gray2.Pix {
					if v != gray1.Pix[i] {
						t.Errorf("got: %v, expected: %v (y=%d, cb=%d, cr=%d)\n",
							gray2.Pix, gray1.Pix, y, cb, cr)
						return
					}
				}
			}
		}
	}
}
