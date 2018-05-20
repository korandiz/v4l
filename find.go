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

import (
	"os"
	"path/filepath"
	"syscall"
)

// FindDevices returns the DeviceInfo for every capture device found in the
// system. The paths are always returned in absolute form.
func FindDevices() []DeviceInfo {
	var (
		devs   = make(map[int]string)
		dInfos []DeviceInfo
	)
	filepath.Walk("/dev", func(path string, fInfo os.FileInfo, _ error) error {
		if fInfo == nil || fInfo.Mode()&os.ModeCharDevice == 0 {
			return nil
		}
		var stat syscall.Stat_t
		if err := syscall.Stat(path, &stat); err != nil {
			return nil
		}
		major, minor := int(stat.Rdev>>8), int(stat.Rdev&255)
		if major != 81 {
			return nil
		}
		if _, ok := devs[minor]; ok {
			return nil
		}
		dev, err := Open(path)
		if err != nil {
			return nil
		}
		defer dev.Close()
		dInfo, err := dev.DeviceInfo()
		if err != nil {
			return nil
		}
		dInfos = append(dInfos, dInfo)
		devs[minor] = path
		return nil
	})
	return dInfos
}
