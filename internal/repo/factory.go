// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zhenorzz/goploy/internal/model"
)

type Repo interface {
	Ping(url string) error
	// Follow the repository code and update to latest
	Follow(projectID int64, target string, url string, branch string) error
	// RemoteBranchList list remote branches in the url
	RemoteBranchList(url string) ([]string, error)
	// BranchList list the local repository's branches
	BranchList(projectID int64) ([]string, error)
	// CommitLog list the local commit log
	CommitLog(projectID int64, rows int) ([]CommitInfo, error)
	// BranchLog list the local commit log from specific branch
	BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error)
	// TagLog list the local commit log from all tag
	TagLog(projectID int64, rows int) ([]CommitInfo, error)
	// CanRollback detect repo can rollback or not
	CanRollback() bool
}

type CommitInfo struct {
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
	Diff      string `json:"diff"`
}

func GetRepo(repoType string) (Repo, error) {
	switch repoType {
	case model.RepoGit:
		return GitRepo{}, nil
	case model.RepoSVN:
		return SvnRepo{}, nil
	case model.RepoFTP:
		return FtpRepo{}, nil
	case model.RepoSFTP:
		return SftpRepo{}, nil
	case model.RepoEmpty:
		return EmptyRepo{}, nil
	}
	return nil, fmt.Errorf("wrong repo type passed")
}

func (commitInfo CommitInfo) ReplaceVars(script string) string {
	scriptVars := map[string]string{
		"${COMMIT_TAG}":       commitInfo.Tag,
		"${COMMIT_BRANCH}":    commitInfo.Branch,
		"${COMMIT_ID}":        commitInfo.Commit,
		"${COMMIT_SHORT_ID}":  commitInfo.Commit,
		"${COMMIT_AUTHOR}":    commitInfo.Author,
		"${COMMIT_TIMESTAMP}": strconv.FormatInt(commitInfo.Timestamp, 10),
		"${COMMIT_MESSAGE}":   commitInfo.Message,
	}

	if len(commitInfo.Commit) > 6 {
		scriptVars["${COMMIT_SHORT_ID}"] = commitInfo.Commit[0:6]
	}

	for key, value := range scriptVars {
		script = strings.Replace(script, key, value, -1)
	}
	return script
}
