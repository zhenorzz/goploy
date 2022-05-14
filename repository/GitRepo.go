// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repository

import (
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"os"
	"strconv"
	"strings"
)

type GitRepo struct{}

func (GitRepo) CanRollback() bool {
	return true
}

// Ping -
func (GitRepo) Ping(url string) error {
	git := utils.GIT{}
	if err := git.LsRemote("-h", url); err != nil {
		return errors.New(git.Err.String())
	}

	return nil
}

// Create -
func (GitRepo) Create(projectID int64) error {
	srcPath := core.GetProjectPath(projectID)
	if _, err := os.Stat(srcPath); err == nil {
		return nil
	}
	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project does not exist, projectID:%d", projectID))
		return err
	}
	if err := os.RemoveAll(srcPath); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to remove, projectID:%d, error: %s", projectID, err.Error()))
		return err
	}
	git := utils.GIT{}
	if err := git.Clone(project.URL, srcPath); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to initialize, projectID:%d, error:%s, detail:%s", projectID, err.Error(), git.Err.String()))
		return err
	}

	git.Dir = srcPath
	if err := git.Current(); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("The project fail to get current branch, projectID:%d, error:%s, detail:%s", projectID, err.Error(), git.Err.String()))
		return err
	}

	currentBranch := utils.ClearNewline(git.Output.String())
	if project.Branch != currentBranch {
		if err := git.Checkout("-b", project.Branch, "origin/"+project.Branch); err != nil {
			core.Log(core.ERROR, fmt.Sprintf("The project fail to switch branch, projectID:%d, error:%s, detail:%s", projectID, err.Error(), git.Err.String()))
			_ = os.RemoveAll(srcPath)
			return err
		}
	}
	core.Log(core.TRACE, fmt.Sprintf("The project success to initialize, projectID:%d", projectID))
	return nil
}

func (gitRepo GitRepo) Follow(project model.Project, target string) error {
	if err := gitRepo.Create(project.ID); err != nil {
		return err
	}
	git := utils.GIT{Dir: core.GetProjectPath(project.ID)}
	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git add .")
	if err := git.Add("."); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}

	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git reset --hard")
	if err := git.Reset("--hard"); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}

	// the length of commit id is 40
	if len(target) != 40 {
		core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git fetch")
		if err := git.Fetch(); err != nil {
			core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
			return errors.New(git.Err.String())
		}
	}

	core.Log(core.TRACE, "projectID:"+strconv.FormatInt(project.ID, 10)+" git checkout -B goploy "+target)
	if err := git.Checkout("-B", "goploy", target); err != nil {
		core.Log(core.ERROR, err.Error()+", detail: "+git.Err.String())
		return errors.New(git.Err.String())
	}
	return nil
}

func (GitRepo) RemoteBranchList(url string) ([]string, error) {
	git := utils.GIT{}
	if err := git.LsRemote("-h", url); err != nil {
		return []string{}, errors.New(git.Err.String())
	}

	var list []string
	for _, branchWithSha := range strings.Split(git.Output.String(), "\n") {
		if len(branchWithSha) != 0 {
			branchWithShaSlice := strings.Fields(branchWithSha)
			branchWithHead := branchWithShaSlice[len(branchWithShaSlice)-1]
			branchWithHeadSlice := strings.Split(branchWithHead, "/")
			list = append(list, branchWithHeadSlice[len(branchWithHeadSlice)-1])
		}
	}
	return list, nil
}

func (GitRepo) BranchList(projectID int64) ([]string, error) {
	git := utils.GIT{Dir: core.GetProjectPath(projectID)}

	if err := git.Fetch(); err != nil {
		return []string{}, errors.New(err.Error() + " detail: " + git.Err.String())
	}

	if err := git.Branch("-r", "--sort=-committerdate"); err != nil {
		return []string{}, errors.New(err.Error() + " detail: " + git.Err.String())
	}

	rawBranchList := strings.Split(git.Output.String(), "\n")

	var list []string
	for _, row := range rawBranchList {
		branch := strings.Trim(row, " ")
		if len(branch) != 0 {
			list = append(list, branch)
		}
	}
	return list, nil
}

func (GitRepo) CommitLog(projectID int64, rows int) ([]CommitInfo, error) {
	git := utils.GIT{Dir: core.GetProjectPath(projectID)}

	if err := git.Log("--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`", "-n", strconv.Itoa(rows)); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	list := parseGITLog(git.Output.String())
	return list, nil
}

func (GitRepo) BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error) {
	git := utils.GIT{Dir: core.GetProjectPath(projectID)}

	if err := git.Log(branch, "--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`", "-n", strconv.Itoa(rows)); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	list := parseGITLog(git.Output.String())
	return list, nil
}

func (GitRepo) TagLog(projectID int64, rows int) ([]CommitInfo, error) {
	git := utils.GIT{Dir: core.GetProjectPath(projectID)}
	if err := git.Add("."); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	if err := git.Reset("--hard"); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	if err := git.Pull(); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	if err := git.Log("--tags", "-n", strconv.Itoa(rows), "--no-walk", "--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`"); err != nil {
		return []CommitInfo{}, errors.New(git.Err.String())
	}

	list := parseGITLog(git.Output.String())
	return list, nil
}

func parseGITLog(rawCommitLog string) []CommitInfo {
	unformatCommitList := strings.Split(rawCommitLog, "`start`")
	unformatCommitList = unformatCommitList[1:]
	var commitList []CommitInfo
	for _, commitRow := range unformatCommitList {
		commitRowSplit := strings.Split(commitRow, "`")
		timestamp, _ := strconv.ParseInt(commitRowSplit[2], 10, 64)
		commitList = append(commitList, CommitInfo{
			Commit:    commitRowSplit[0],
			Author:    commitRowSplit[1],
			Timestamp: timestamp,
			Message:   commitRowSplit[3],
			Tag:       commitRowSplit[4],
			Diff:      strings.Trim(commitRowSplit[5], "\n"),
		})
	}
	return commitList
}
