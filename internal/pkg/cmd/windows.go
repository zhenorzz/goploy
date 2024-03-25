// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
)

type WindowsCmd struct{}

func (w WindowsCmd) Script(mode, file string) string {
	if mode == "" || mode == "cmd" {
		mode = "cmd /C"
	}
	return fmt.Sprintf("%s %s", mode, file)
}

func (w WindowsCmd) ChangeDirTime(dir string) string {
	tmpFile := Join(dir, "goploy.tmp")
	return fmt.Sprintf("type nul > %s && del %[1]s", tmpFile)
}

func (w WindowsCmd) Symlink(src, target string) string {
	return fmt.Sprintf("mklink /D %s %s", target, src)
}

func (w WindowsCmd) Remove(file string) string {
	return fmt.Sprintf("del /Q %s", file)
}
