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

import (
	"io"
	"math/rand"
	"testing"
)

const N = 100

func TestBuffer_Size(t *testing.T) {
	buf := initBufferTest()
	for i := 0; i < N+10; i++ {
		size := buf.Size()
		if size != N {
			t.Errorf("got: %d, expected: %d (i=%d)\n", size, N, i)
			return
		}
		buf.ReadByte()
	}
}

func TestBuffer_Len(t *testing.T) {
	buf := initBufferTest()
	for i := 0; i < N+10; i++ {
		length := buf.Len()
		n := N - i
		if n < 0 {
			n = 0
		}
		if length != n {
			t.Errorf("got: %d, expected: %d (i=%d)\n", length, n, i)
			return
		}
		buf.ReadByte()
	}
}

func TestBuffer_Read(t *testing.T) {
	var (
		buf  = initBufferTest()
		arr  = make([]byte, 2*N)
		err  error
		n, m int
	)
	for err == nil {
		n, err = buf.Read(arr[m : m+15])
		m += n
	}
	if err != io.EOF {
		t.Errorf("expected EOF, got: %v (m=%d)\n", err, m)
		return
	}
	if m != N {
		t.Errorf("wrong byte count: %d (expected: %d)\n", m, N)
		return
	}
	for i := 0; i < N; i++ {
		if arr[i] != byte(i) {
			t.Errorf("got: %d, expected: %d (i=%d)\n", arr[i], byte(i), i)
			return
		}
	}
}

func TestBuffer_ReadAt(t *testing.T) {
	buf := initBufferTest()
	arr := make([]byte, N+10)
	for start := -10; start < N+10; start++ {
		for length := 0; length < N+10; length++ {
			n, err := buf.ReadAt(arr[:length], int64(start))
			if start < 0 {
				if err == nil || n != 0 {
					t.Errorf("read at %d should have failed (length=%d, n=%d)\n",
						start, length, n)
					return
				}
				continue
			}
			if n > length {
				t.Errorf("read too many bytes (start=%d, length=%d, n=%d)\n",
					start, length, n)
				return
			}
			if n < length && err == nil {
				t.Errorf("short read with nil error (start=%d, length=%d, n=%d)\n",
					start, length, n)
				return
			}
			if start+length <= N && n < length {
				t.Errorf("short read (start=%d, length=%d, n=%d)\n",
					start, length, n)
				return
			}
			m := N - start
			if m < 0 {
				m = 0
			}
			if start+length > N && n != m {
				t.Errorf("read %d bytes instead of %d (start=%d, length=%d)\n",
					n, m, start, length)
				return
			}
			if start+length > N && err != io.EOF {
				t.Errorf("expected EOF, got: %v (start=%d, length=%d)\n",
					err, start, length)
				return
			}
			for i := 0; i < n; i++ {
				if arr[i] != byte(start+i) {
					t.Errorf("expected: %d, got: %d (start=%d, length=%d, i=%d)\n",
						start+i, arr[i], start, length, i)
					return
				}
			}
		}
	}
}

func TestBuffer_ReadByte(t *testing.T) {
	buf := initBufferTest()
	for i := 0; i < N; i++ {
		b, err := buf.ReadByte()
		if err != nil {
			t.Errorf("unexpected error: %v (i=%d)\n", err, i)
			return
		}
		if b != byte(i) {
			t.Errorf("got: %d, expected: %d\n", b, i)
			return
		}
	}
	if _, err := buf.ReadByte(); err != io.EOF {
		t.Error("expected EOF, got:", err)
		return
	}
}

func TestBuffer_Seek(t *testing.T) {
	buf := initBufferTest()
	rand.Seed(42)
	for _, whence := range []int{io.SeekStart, io.SeekEnd} {
		for i := 0; i < 10*N; i++ {
			offs := rand.Int63n(N)
			var m int64
			if whence == io.SeekEnd {
				offs = -offs - 1
				m = N
			}
			m += offs
			n, err := buf.Seek(offs, whence)
			if err != nil {
				t.Errorf("unexpected error: %v (i=%d, offs=%d, whence=%d)\n",
					err, i, offs, whence)
				return
			}
			if n != m {
				t.Errorf("bad offset: %d (m=%d, i=%d, offs=%d, whence=%d)\n",
					n, m, i, offs, whence)
				return
			}
			b, err := buf.ReadByte()
			if err != nil {
				t.Errorf("unexpected error: %v (i=%d, offs=%d, whence=%d)\n",
					err, i, offs, whence)
				return
			}
			if b != byte(m) {
				t.Errorf("read: %d (m=%d, i=%d, offs=%d, whence=%d)\n",
					b, m, i, offs, whence)
				return
			}
		}
	}
	buf = initBufferTest()
	var p int64
	for i := 0; i < 10*N; i++ {
		m := rand.Int63n(N)
		offs := m - p
		n, err := buf.Seek(offs, io.SeekCurrent)
		if err != nil {
			t.Errorf("unexpected error: %v (i=%d, offs=%d, whence=%d)\n",
				err, i, offs, io.SeekCurrent)
			return
		}
		if n != m {
			t.Errorf("bad offset: %d (m=%d, i=%d, offs=%d, whence=%d)\n",
				n, m, i, offs, io.SeekCurrent)
			return
		}
		b, err := buf.ReadByte()
		if err != nil {
			t.Errorf("unexpected error: %v (i=%d, offs=%d, whence=%d)\n",
				err, i, offs, io.SeekCurrent)
			return
		}
		if b != byte(m) {
			t.Errorf("read: %d (m=%d, i=%d, offs=%d, whence=%d)\n",
				b, m, i, offs, io.SeekCurrent)
			return
		}
		p = n + 1
	}
	_, err := buf.Seek(-1, io.SeekStart)
	if err == nil {
		t.Error("seek to offset -1 should have failed")
		return
	}
}

func initBufferTest() *Buffer {
	buf := make([]byte, N)
	for i := range buf {
		buf[i] = byte(i)
	}
	d := device{
		buffers:   [][]byte{nil, nil, buf, nil},
		bufIndex:  2,
		nCaptures: 1,
	}
	return &Buffer{
		d:   &d,
		n:   1,
		pos: 0,
	}
}
