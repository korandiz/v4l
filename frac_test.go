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

import "testing"

func TestFrac_Reduce(t *testing.T) {
	var x = []struct{ in, out Frac }{
		{Frac{0, 1}, Frac{0, 1}},
		{Frac{0, 2}, Frac{0, 1}},
		{Frac{0, 3}, Frac{0, 1}},
		{Frac{0, 4}, Frac{0, 1}},
		{Frac{0, 6}, Frac{0, 1}},
		{Frac{0, 12}, Frac{0, 1}},
		{Frac{1, 0}, Frac{1, 0}},
		{Frac{1, 1}, Frac{1, 1}},
		{Frac{1, 2}, Frac{1, 2}},
		{Frac{1, 3}, Frac{1, 3}},
		{Frac{1, 4}, Frac{1, 4}},
		{Frac{1, 6}, Frac{1, 6}},
		{Frac{1, 12}, Frac{1, 12}},
		{Frac{2, 0}, Frac{1, 0}},
		{Frac{2, 1}, Frac{2, 1}},
		{Frac{2, 2}, Frac{1, 1}},
		{Frac{2, 3}, Frac{2, 3}},
		{Frac{2, 4}, Frac{1, 2}},
		{Frac{2, 6}, Frac{1, 3}},
		{Frac{2, 12}, Frac{1, 6}},
		{Frac{3, 0}, Frac{1, 0}},
		{Frac{3, 1}, Frac{3, 1}},
		{Frac{3, 2}, Frac{3, 2}},
		{Frac{3, 3}, Frac{1, 1}},
		{Frac{3, 4}, Frac{3, 4}},
		{Frac{3, 6}, Frac{1, 2}},
		{Frac{3, 12}, Frac{1, 4}},
		{Frac{4, 0}, Frac{1, 0}},
		{Frac{4, 1}, Frac{4, 1}},
		{Frac{4, 2}, Frac{2, 1}},
		{Frac{4, 3}, Frac{4, 3}},
		{Frac{4, 4}, Frac{1, 1}},
		{Frac{4, 6}, Frac{2, 3}},
		{Frac{4, 12}, Frac{1, 3}},
		{Frac{6, 0}, Frac{1, 0}},
		{Frac{6, 1}, Frac{6, 1}},
		{Frac{6, 2}, Frac{3, 1}},
		{Frac{6, 3}, Frac{2, 1}},
		{Frac{6, 4}, Frac{3, 2}},
		{Frac{6, 6}, Frac{1, 1}},
		{Frac{6, 12}, Frac{1, 2}},
		{Frac{12, 0}, Frac{1, 0}},
		{Frac{12, 1}, Frac{12, 1}},
		{Frac{12, 2}, Frac{6, 1}},
		{Frac{12, 3}, Frac{4, 1}},
		{Frac{12, 4}, Frac{3, 1}},
		{Frac{12, 6}, Frac{2, 1}},
		{Frac{12, 12}, Frac{1, 1}},
		{Frac{54665710, 239322993}, Frac{54665710, 239322993}},
		{Frac{239322993, 54665710}, Frac{239322993, 54665710}},
		{Frac{223092870, 1078282205}, Frac{6, 29}},
		{Frac{1078282205, 223092870}, Frac{29, 6}},
		{Frac{0, 0}, Frac{0, 0}},
	}

	for _, y := range x {
		out := y.in.Reduce()
		if out != y.out {
			t.Errorf("in: %d/%d, out: %d/%d, expected: %d/%d\n",
				y.in.N, y.in.D, out.N, out.D, y.out.N, y.out.D)
		}
	}
}

func TestFrac_Cmp(t *testing.T) {
	rel := map[int]string{-1: "<", 0: "=", 1: ">"}
	var x = [][]Frac{
		{Frac{0, 1}, Frac{0, 2}, Frac{0, 3}, Frac{0, 4}, Frac{0, 6}, Frac{0, 12}},
		{Frac{1, 12}},
		{Frac{1, 6}, Frac{2, 12}},
		{Frac{223092870, 1078282205}, Frac{6, 29}},
		{Frac{54665710, 239322993}},
		{Frac{1, 4}, Frac{3, 12}},
		{Frac{1, 3}, Frac{2, 6}, Frac{4, 12}},
		{Frac{1, 2}, Frac{2, 4}, Frac{3, 6}, Frac{6, 12}},
		{Frac{2, 3}, Frac{4, 6}},
		{Frac{3, 4}},
		{Frac{1, 1}, Frac{2, 2}, Frac{3, 3}, Frac{4, 4}, Frac{6, 6}, Frac{12, 12}},
		{Frac{4, 3}},
		{Frac{3, 2}, Frac{6, 4}},
		{Frac{2, 1}, Frac{4, 2}, Frac{6, 3}, Frac{12, 6}},
		{Frac{3, 1}, Frac{6, 2}, Frac{12, 4}},
		{Frac{4, 1}, Frac{12, 3}},
		{Frac{239322993, 54665710}},
		{Frac{1078282205, 223092870}},
		{Frac{6, 1}, Frac{12, 2}},
		{Frac{12, 1}},
		{Frac{1, 0}, Frac{2, 0}, Frac{3, 0}, Frac{4, 0}, Frac{6, 0}, Frac{12, 0}},
	}

	for i1 := range x {
		for _, f := range x[i1] {
			for i2 := range x {
				for _, g := range x[i2] {
					var r int
					if i1 < i2 {
						r = -1
					} else if i1 > i2 {
						r = 1
					}
					cmp := f.Cmp(g)
					if cmp != r {
						t.Errorf("got: %d/%d %s %d/%d, expected: %d/%d %s %d/%d\n",
							f.N, f.D, rel[cmp], g.N, g.D,
							f.N, f.D, rel[r], g.N, g.D)

					}
				}
			}
		}
	}
}
