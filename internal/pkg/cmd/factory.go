// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package cmd

import (
	"strings"
)

type Cmd interface {
	Symlink(src, target string) string
	Remove(file string) string
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

func ExtractSeparator(path string) byte {
	if strings.Contains(path, "\\") {
		return '\\'
	}
	return '/'
}

func Join(elem ...string) string {
	size := 0
	for _, e := range elem {
		size += len(e)
	}
	if size == 0 {
		return ""
	}
	sep := ExtractSeparator(elem[0])
	buf := make([]byte, 0, size+len(elem)-1)
	for _, e := range elem {
		if len(buf) > 0 || e != "" {
			if len(buf) > 0 {
				buf = append(buf, sep)
			}
			buf = append(buf, e...)
		}
	}
	return string(buf)
}
