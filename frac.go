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

// A Frac represents the fractional number N/D.
type Frac struct {
	N uint32
	D uint32
}

// Cmp returns -1, 0, or 1 if f is less then, equal to, or greater than g,
// respectively.
func (f Frac) Cmp(g Frac) int {
	F := uint64(f.N) * uint64(g.D)
	G := uint64(g.N) * uint64(f.D)
	switch {
	case F < G:
		return -1
	case F > G:
		return 1
	default:
		return 0
	}
}

// Reduce returns f reduced to lowest terms.
func (f Frac) Reduce() Frac {
	gcd, r := f.N, f.D
	for r != 0 {
		gcd, r = r, gcd%r
	}
	if gcd == 0 {
		gcd = 1
	}
	return Frac{f.N / gcd, f.D / gcd}
}
