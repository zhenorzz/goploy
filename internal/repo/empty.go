// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repo

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
)

type EmptyRepo struct{}

func (EmptyRepo) CanRollback() bool {
	return false
}

func (repo EmptyRepo) Ping(url string) error {
	return nil
}

func (repo EmptyRepo) Follow(projectID int64, _ string, projectURL string, _ string) error {
	srcPath := config.GetProjectPath(projectID)
	_ = os.RemoveAll(srcPath)
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		log.Error(fmt.Sprintf("The project fail to mkdir, projectID:%d, error:%s", projectID, err.Error()))
		return err
	}
	return nil
}

func (EmptyRepo) RemoteBranchList(url string) ([]string, error) {
	return []string{"virtual"}, nil
}

func (EmptyRepo) BranchList(projectID int64) ([]string, error) {
	return []string{"virtual"}, nil
}

func (EmptyRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
	commitInfo := CommitInfo{
		Branch:    "virtual",
		Commit:    "",
		Author:    "",
		Timestamp: time.Now().Unix(),
		Message:   "",
		Tag:       "",
		Diff:      "",
	}
	return []CommitInfo{commitInfo}, nil
}

func (repo EmptyRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	return []CommitInfo{{Commit: "virtual"}}, nil
}

func (repo EmptyRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	return []CommitInfo{}, nil
}
