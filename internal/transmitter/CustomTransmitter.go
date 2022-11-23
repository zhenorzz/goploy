// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package transmitter

import (
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/model"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type customTransmitter struct {
	Project       model.Project
	ProjectServer model.ProjectServer
}

func (ct customTransmitter) String() string {
	project := ct.Project
	server := ct.ProjectServer
	script := project.TransferOption
	scriptVars := map[string]string{
		"${PROJECT_ID}":            strconv.FormatInt(project.ID, 10),
		"${PROJECT_PATH}":          project.Path,
		"${PROJECT_SYMLINK_PATH}":  path.Join(project.SymlinkPath, project.LastPublishToken),
		"${PROJECT_NAME}":          project.Name,
		"${PROJECT_BRANCH}":        project.Branch,
		"${REPOSITORY_TYPE}":       project.RepoType,
		"${REPOSITORY_URL}":        project.URL,
		"${REPOSITORY_PATH}":       config.GetProjectPath(project.ID),
		"${AFTER_DEPLOY_FILENAME}": fmt.Sprintf("goploy-after-deploy-p%d-s%d.%s", project.ID, server.ServerID, pkg.GetScriptExt(project.AfterDeployScriptMode)),
		"${PUBLISH_TOKEN}":         project.LastPublishToken,
		"${SERVER_ID}":             strconv.FormatInt(server.ServerID, 10),
		"${SERVER_NAME}":           server.ServerName,
		"${SERVER_IP}":             server.ServerIP,
		"${SERVER_PORT}":           strconv.Itoa(server.ServerPort),
		"${SERVER_OWNER}":          server.ServerOwner,
		"${SERVER_PASSWORD}":       server.ServerPassword,
		"${SERVER_PATH}":           server.ServerPath,
		"${SERVER_JUMP_IP}":        server.ServerJumpIP,
		"${SERVER_JUMP_PORT}":      strconv.Itoa(server.ServerJumpPort),
		"${SERVER_JUMP_OWNER}":     server.ServerJumpOwner,
		"${SERVER_JUMP_PASSWORD}":  server.ServerJumpPassword,
		"${SERVER_JUMP_PATH}":      server.ServerJumpPath,
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
