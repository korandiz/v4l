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

import "io"

// A Buffer holds the raw image data of a frame captured from a Device. It
// implements io.Reader, io.ByteReader, io.ReaderAt, and io.Seeker. A call to
// Capture, Close, or TurnOff on the corresponding Device may cause the contents
// of the buffer to go away.
type Buffer struct {
	d   *device
	n   uint64
	pos int
	seq uint32
}

// Size returns the total number of bytes in the buffer. As long as the data is
// available, the return value is constant and unaffected by calls to the
// methods of Buffer. If the data is no longer available, it returns 0.
func (b *Buffer) Size() int64 {
	return int64(len(b.source()))
}

// Len returns the number of unread bytes in the buffer. If the data is no
// longer available, it returns 0.
func (b *Buffer) Len() int {
	src := b.source()
	if src == nil {
		return 0
	}
	return len(src) - b.pos
}

// SeqNum returns the sequence number of the frame in the buffer as reported by
// the kernel.
func (b *Buffer) SeqNum() uint32 {
	return b.seq
}

// Read reads up to len(dst) bytes into dst, and returns the number of bytes
// read, along with any error encountered.
func (b *Buffer) Read(dst []byte) (int, error) {
	src := b.source()
	if src == nil {
		return 0, ErrBufferGone
	}
	if b.pos == len(src) {
		return 0, io.EOF
	}
	n := copy(dst, src[b.pos:])
	b.pos += n
	return n, nil
}

// ReadAt reads up to len(dst) bytes into dst starting at the specified offset,
// and returns the number of bytes read, along with any error encountered.
// The seek offset is unaffected by ReadAt.
func (b *Buffer) ReadAt(dst []byte, offset int64) (int, error) {
	src := b.source()
	if src == nil {
		return 0, ErrBufferGone
	}
	if offset < 0 {
		return 0, Error("negative offset")
	}
	if offset >= int64(len(src)) {
		return 0, io.EOF
	}
	n := copy(dst, src[offset:])
	if n < len(dst) {
		return n, io.EOF
	}
	return n, nil
}

// ReadByte returns the next byte in the buffer.
func (b *Buffer) ReadByte() (byte, error) {
	src := b.source()
	if src == nil {
		return 0, ErrBufferGone
	}
	if b.pos == len(src) {
		return 0, io.EOF
	}
	x := src[b.pos]
	b.pos++
	return x, nil
}

// Seek sets the seek offset.
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	src := b.source()
	if src == nil {
		return 0, ErrBufferGone
	}
	var i int64
	switch whence {
	case io.SeekStart:
		i = offset
	case io.SeekCurrent:
		i = int64(b.pos) + offset
	case io.SeekEnd:
		i = int64(len(src)) + offset
	default:
		return 0, Error("invalid whence")
	}
	if i < 0 {
		return 0, Error("negative offset")
	}
	if i > int64(len(src)) {
		i = int64(len(src))
	}
	b.pos = int(i)
	return i, nil
}

// source returns the underlying byte slice of the buffer, or nil, if it's no
// longer available.
func (b *Buffer) source() []byte {
	if b.d.nCaptures != b.n || b.d.bufIndex == noBuffer {
		return nil
	}
	return b.d.buffers[b.d.bufIndex]
}
