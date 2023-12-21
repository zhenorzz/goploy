// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package transmitter

import (
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type sftpTransmitter struct {
	Project       model.Project
	ProjectServer model.ProjectServer
}

func (st sftpTransmitter) String() string {
	return "sftp " + st.Project.ReplaceVars(st.Project.TransferOption)
}

func (st sftpTransmitter) Exec() (string, error) {
	client, err := st.ProjectServer.ToSSHConfig().Dial()
	if err != nil {
		return "", err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()

	project := st.Project
	transferOption, _ := pkg.ParseCommandLine(st.ProjectServer.ReplaceVars(project.ReplaceVars(project.TransferOption)))
	var opt struct {
		isVerbose bool
		isDelete  bool
	}
	var excludes []string
	var excludeRegexps []*regexp.Regexp
	var includes []string
	var includeRegexps []*regexp.Regexp
	output := "sending file list\n"
	nextItem := ""
	for _, item := range transferOption {
		if item == "--exclude" {
			nextItem = "--exclude"
			continue
		} else if item == "--exclude-regexp" {
			nextItem = "--exclude-regexp"
			continue
		} else if item == "--include" {
			nextItem = "--include"
			continue
		} else if item == "--include-regexp" {
			nextItem = "--include-regexp"
			continue
		} else if strings.HasPrefix(item, "-") {
			if strings.Contains(item, "v") {
				opt.isVerbose = true
			}
		} else if item == "--delete" {
			opt.isDelete = true
		}

		if nextItem == "--exclude" {
			excludes = append(excludes, item)
		} else if nextItem == "--exclude-regexp" {
			r, err := regexp.Compile(item)
			if err != nil {
				return "", err
			}
			excludeRegexps = append(excludeRegexps, r)
		} else if nextItem == "--include" {
			tempItem := ""
			for _, s := range strings.Split(item, "/") {
				tempItem = tempItem + s
				includes = append(includes, tempItem)
				tempItem = tempItem + "/"
			}
			includes = append(includes, item)
		} else if nextItem == "--include-regexp" {
			r, err := regexp.Compile(item)
			if err != nil {
				return "", err
			}
			includeRegexps = append(includeRegexps, r)
		}
		nextItem = ""
	}
	includes = append(includes, project.Script.AfterDeploy.ScriptNames...)

	srcPath := config.GetProjectPath(project.ID) + "/"
	destPath := project.Path
	if len(project.SymlinkPath) != 0 {
		destPath = path.Join(project.SymlinkPath, project.LastPublishToken)
	}
	var uploadViaSFTP func(localDir, remoteDir string) error
	uploadViaSFTP = func(localDir, remoteDir string) error {
		localEntries, err := os.ReadDir(localDir)

		if err != nil {
			return err
		}
		for _, entry := range localEntries {
			nextLocalDir := path.Join(localDir, entry.Name())
			nextRemoteDir := path.Join(remoteDir, entry.Name())
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			target := nextLocalDir[len(srcPath):]

			isExclude := false
			for _, exclude := range excludes {
				if target == exclude {
					isExclude = true
					break
				}
			}

			for _, excludeRegexp := range excludeRegexps {
				if excludeRegexp.MatchString(target) {
					isExclude = true
					break
				}
			}

			if isExclude {
				for _, include := range includes {
					if target == include {
						isExclude = false
						break
					}
				}

				for _, includeRegexp := range includeRegexps {
					if includeRegexp.MatchString(target) {
						isExclude = false
						break
					}
				}
			}

			if isExclude {
				goto endSync
			} else {
				goto startSync
			}

		startSync:
			output += target + "\n"
			if entry.IsDir() {
				if err = sftpClient.MkdirAll(nextRemoteDir); err != nil {
					return err
				}
				if err = uploadViaSFTP(nextLocalDir, nextRemoteDir); err != nil {
					return err
				}
			} else {
				remoteFile, err := sftpClient.Create(nextRemoteDir)
				if err != nil {
					return err
				}

				localFile, err := os.Open(nextLocalDir)
				if err != nil {
					remoteFile.Close()
					return err
				}
				_, err = io.Copy(remoteFile, localFile)
				localFile.Close()
				remoteFile.Close()
				if err != nil {
					println(nextLocalDir, nextRemoteDir)
					return err
				}
			}
		endSync:
		}
		return nil
	}

	if opt.isDelete {
		if err := deleteDirectory(sftpClient, destPath); err != nil {
			return "", err
		}
	}

	if err := sftpClient.MkdirAll(destPath); err != nil {
		return "", err
	}

	if err := uploadViaSFTP(srcPath, destPath); err != nil {
		return "", err
	}

	if opt.isVerbose {
		output += "sent files success"
	} else {
		output = "sent files success"
	}
	return output, nil
}

func deleteDirectory(sftpClient *sftp.Client, dirPath string) error {
	fileInfos, err := sftpClient.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		filePath := dirPath + "/" + fileInfo.Name()

		if fileInfo.IsDir() {
			err = deleteDirectory(sftpClient, filePath)
			if err != nil {
				return err
			}
		} else {
			err = sftpClient.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
