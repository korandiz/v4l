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

// Package yuyv provides support for the YUYV format.
package yuyv

import (
	"image"
	"image/color"
)

// FourCC of the YUYV format.
const FourCC = 'Y' | 'U'<<8 | 'Y'<<16 | 'V'<<24

// An Image is an in-memory YUYV image. It implements the image.Image interface.
type Image struct {
	// Pix holds the pixel data. Every four bytes corresponds to a pair of
	// horizontally adjacent pixels. The two pixels have separate luma values
	// and shared chroma. The sample order is Y0, Cb, Y1, Cr.
	Pix []uint8

	// Stride is the distance in bytes between vertically adjacent pixels.
	// Line y starts at index (y-Rect.Min.Y)*Stride.
	Stride int

	// Rect is the bounds of the image.
	Rect image.Rectangle
}

// New returns a new YUYV image with the given bounds.
func New(rect image.Rectangle) *Image {
	l, r := rect.Min.X&^1, (rect.Max.X+1)&^1
	w, h := r-l, rect.Dy()
	return &Image{make([]uint8, 2*w*h), 2 * w, rect}
}

// ColorModel returns color.YCbCrModel.
func (img *Image) ColorModel() color.Model {
	return color.YCbCrModel
}

// Bounds returns the bounds of the image.
func (img *Image) Bounds() image.Rectangle {
	return img.Rect
}

// At returns the color of the pixel at (x, y).
func (img *Image) At(x, y int) color.Color {
	return img.YCbCrAt(x, y)
}

// YCbCrAt returns the color of the pixel at (x, y).
func (img *Image) YCbCrAt(x, y int) color.YCbCr {
	if !(image.Point{x, y}.In(img.Rect)) {
		return color.YCbCr{}
	}
	return color.YCbCr{
		img.Pix[img.YOffset(x, y)],
		img.Pix[img.CbOffset(x, y)],
		img.Pix[img.CrOffset(x, y)],
	}
}

// YOffset returns the index at which the Y component of the pixel at (x, y) is
// located in Pix.
func (img *Image) YOffset(x, y int) int {
	return img.PixPairOffset(x, y) + 2*(x&1)
}

// CbOffset returns the index at which the Cb component of the pixel at (x, y)
// is located in Pix.
func (img *Image) CbOffset(x, y int) int {
	return img.PixPairOffset(x, y) + 1
}

// CrOffset returns the index at which the Cr component of the pixel at (x, y)
// is located in Pix.
func (img *Image) CrOffset(x, y int) int {
	return img.PixPairOffset(x, y) + 3
}

// PixPairOffset returns the index of the first element in Pix that corresponds
// to the pixel pair to which (x, y) belongs.
func (img *Image) PixPairOffset(x, y int) int {
	return (y-img.Rect.Min.Y)*img.Stride + (x&^1-img.Rect.Min.X&^1)*2
}

// SubImage returns an image representing the portion of img visible through
// rect. The returned value shares pixels with the original image.
func (img *Image) SubImage(rect image.Rectangle) image.Image {
	rect = rect.Intersect(img.Rect)
	if rect.Empty() {
		return &Image{}
	}
	i := img.PixPairOffset(rect.Min.X, rect.Min.Y)
	return &Image{img.Pix[i:], img.Stride, rect}
}

// Opaque returns true.
func (img *Image) Opaque() bool {
	return true
}
