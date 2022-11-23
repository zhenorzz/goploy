// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"
)

type WindowsCmd struct{}

func (w WindowsCmd) Script(mode, file string) string {
	if mode == "" || mode == "cmd" {
		mode = "cmd /C"
	}
	return fmt.Sprintf("%s %s", mode, w.Path(file))
}

func (w WindowsCmd) ChangeDirTime(dir string) string {
	tmpFile := w.Path(dir + "/goploy.tmp")
	return fmt.Sprintf("type nul > %s && del %[1]s", tmpFile)
}

func (WindowsCmd) Path(file string) string {
	return strings.ReplaceAll(file, "/", "\\")
}

func (w WindowsCmd) Symlink(src, target string) string {
	return fmt.Sprintf("mklink /D %s %s", w.Path(target), w.Path(src))
}

func (w WindowsCmd) Remove(file string) string {
	return fmt.Sprintf("del /Q %s", w.Path(file))
}
