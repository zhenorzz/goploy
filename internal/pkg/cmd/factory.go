// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package cmd

type Cmd interface {
	Symlink(src, target string) string
	Remove(file string) string
	Path(file string) string
	ChangeDirTime(dir string) string
	Script(mode, file string) string
}

func New(os string) Cmd {
	if os == "windows" {
		return WindowsCmd{}
	} else {
		return LinuxCmd{}
	}
}
