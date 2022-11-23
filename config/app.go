// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	AssetDir string
)

func GetAssetDir() string {
	if AssetDir != "" {
		return AssetDir
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	app, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(app, "/")
	if i < 0 {
		i = strings.LastIndex(app, "\\")
	}
	if i < 0 {
		panic(err)
	}
	AssetDir = app[0 : i+1]
	return AssetDir
}

func GetConfigFile() string {
	return path.Join(GetAssetDir(), "goploy.toml")
}

func GetPidFile() string {
	return path.Join(GetAssetDir(), "goploy.pid")
}

func GetRepositoryPath() string {
	if Toml.APP.RepositoryPath != "" {
		return path.Join(Toml.APP.RepositoryPath, "repository")
	}
	return path.Join(GetAssetDir(), "repository")
}

func GetProjectFilePath(projectID int64) string {
	return path.Join(GetRepositoryPath(), "project-file", "project_"+strconv.FormatInt(projectID, 10))
}

func GetProjectPath(projectID int64) string {
	return path.Join(GetRepositoryPath(), "project_"+strconv.FormatInt(projectID, 10))
}

func GetTerminalLogPath(tlID int64) string {
	return path.Join(GetRepositoryPath(), "terminal-log", strconv.FormatInt(tlID, 10)+".cast")
}
