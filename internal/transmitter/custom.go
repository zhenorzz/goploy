// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package transmitter

import (
	"fmt"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"os/exec"
	"strings"
)

type customTransmitter struct {
	Project       model.Project
	ProjectServer model.ProjectServer
}

func (ct customTransmitter) String() string {
	project := ct.Project
	server := ct.ProjectServer
	script := server.ReplaceVars(project.ReplaceVars(project.TransferOption))
	scriptVars := map[string]string{
		"${AFTER_DEPLOY_FILENAME}": fmt.Sprintf("goploy-after-deploy-p%d-s%d.%s", project.ID, server.ServerID, pkg.GetScriptExt(project.Script.AfterDeploy.Mode)),
	}
	for index, scriptName := range project.Script.AfterDeploy.ScriptNames {
		scriptVars[fmt.Sprintf("${AFTER_DEPLOY_FILENAME_YAML_%d}", index)] = scriptName
	}
	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}

func (ct customTransmitter) Exec() (string, error) {
	if ct.Project.TransferOption == "" {
		return "", nil
	}
	parts, _ := pkg.ParseCommandLine(ct.String())
	cmd := exec.Command(parts[0], parts[1:]...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return string(output), err
	} else {
		return string(output), nil
	}
}
