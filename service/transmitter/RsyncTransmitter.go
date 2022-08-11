// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package transmitter

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type rsyncTransmitter struct {
	Project       model.Project
	ProjectServer model.ProjectServer
}

func (rt rsyncTransmitter) Args() []string {
	project := rt.Project
	projectServer := rt.ProjectServer
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destDir := project.Path
	if len(project.SymlinkPath) != 0 {
		destDir = path.Join(project.SymlinkPath, project.LastPublishToken)
	}
	srcPath := core.GetProjectPath(project.ID) + "/"
	// rsync path can not contain colon
	// windows like C:\
	if strings.Contains(srcPath, ":\\") {
		srcPath = "/cygdrive/" + strings.Replace(srcPath, ":\\", "/", 1)
	}

	rsyncOption, _ := utils.ParseCommandLine(project.TransferOption)
	rsyncOption = append([]string{
		"--include",
		"goploy-after-deploy." + utils.GetScriptExt(project.AfterDeployScriptMode),
	}, rsyncOption...)
	rsyncOption = append([]string{
		"--exclude",
		"goploy-after-pull." + utils.GetScriptExt(project.AfterPullScriptMode),
	}, rsyncOption...)
	rsyncOption = append(rsyncOption, "-e", projectServer.ToSSHOption())

	destPath := remoteMachine + ":" + destDir
	rsyncOption = append(rsyncOption, "--rsync-path=mkdir -p "+destDir+" && rsync", srcPath, destPath)
	return rsyncOption
}

func (rt rsyncTransmitter) String() string {
	logRsyncCmd := regexp.MustCompile(`sshpass -p .*\s`).
		ReplaceAllString(exec.Command("rsync", rt.Args()...).String(), "sshpass -p ***** ")
	return logRsyncCmd
}

func (rt rsyncTransmitter) Exec() (string, error) {
	// example
	// rsync -rtv -e "ssh -o StrictHostKeyChecking=no -p 22 -i C:\Users\Administrator\.ssh\id_rsa" --rsync-path="mkdir -p /data/www/test && rsync" ./main.go root@127.0.0.1:/tmp/test/
	cmd := exec.Command("rsync", rt.Args()...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return string(output), err
	} else {
		return string(output), nil
	}
}
