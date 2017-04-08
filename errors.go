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

// +build linux

package v4l

// An Error is simply an error message.
type Error string

// Error returns e as a string.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrWrongDevice is returned by Open when attempting to open a file that is
	// not a V4L capture device.
	ErrWrongDevice = Error("not a V4L capture device")

	// ErrUnsupported indicates that an operation failed due to a limitation of
	// this library.
	ErrUnsupported = Error("unsupported device or operation")

	// ErrBufferGone is returned by methods of Buffer when the contents of the
	// buffer is no longer available.
	ErrBufferGone = Error("buffer contents not available")
)
