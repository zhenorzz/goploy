// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package transmitter

import (
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type rsyncTransmitter struct {
	Project       model.Project
	ProjectServer model.ProjectServer
}

func (rt rsyncTransmitter) args() []string {
	project := rt.Project
	projectServer := rt.ProjectServer
	remoteMachine := projectServer.ServerOwner + "@" + projectServer.ServerIP
	destDir := project.Path
	if len(project.SymlinkPath) != 0 {
		destDir = path.Join(project.SymlinkPath, project.LastPublishToken)
	}
	srcPath := config.GetProjectPath(project.ID) + "/"
	// rsync path can not contain colon
	// windows like C:\
	if strings.Contains(srcPath, ":\\") {
		srcPath = "/cygdrive/" + strings.Replace(srcPath, ":\\", "/", 1)
	}

	rsyncOption, _ := pkg.ParseCommandLine(project.ReplaceVars(project.TransferOption))
	rsyncOption = append([]string{
		"--include",
		fmt.Sprintf("goploy-after-deploy-p%d-s%d.%s", project.ID, projectServer.ServerID, pkg.GetScriptExt(project.Script.AfterDeploy.Mode)),
	}, rsyncOption...)
	rsyncOption = append(rsyncOption, "-e", projectServer.ToSSHOption())

	if projectServer.ServerOS == model.ServerOSLinux {
		rsyncOption = append(rsyncOption, "--rsync-path=mkdir -p "+destDir+" && rsync")
	}

	destPath := remoteMachine + ":" + destDir
	rsyncOption = append(rsyncOption, srcPath, destPath)
	return rsyncOption
}

func (rt rsyncTransmitter) String() string {
	logRsyncCmd := regexp.MustCompile(`sshpass -p .*\s`).
		ReplaceAllString(exec.Command("rsync", rt.args()...).String(), "sshpass -p ***** ")
	return logRsyncCmd
}

func (rt rsyncTransmitter) Exec() (string, error) {
	// example
	// rsync -rtv -e "ssh -o StrictHostKeyChecking=no -p 22 -i C:\Users\Administrator\.ssh\id_rsa" --rsync-path="mkdir -p /data/www/test && rsync" ./main.go root@127.0.0.1:/tmp/test/
	cmd := exec.Command("rsync", rt.args()...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return string(output), err
	} else {
		return string(output), nil
	}
}
