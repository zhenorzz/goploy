package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"strconv"
	"strings"
)

// Repository struct
type Repository Controller

// GetCommitList get latest 10 commit list
func (Repository) GetCommitList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.GetProjectPath(project.ID)
	git := utils.GIT{Dir: srcPath}
	if err := git.Log(gp.URLQuery.Get("branch"), "--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`", "-n", "10"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	commitList := utils.ParseGITLog(git.Output.String())

	return &core.Response{
		Data: struct {
			CommitList []utils.Commit `json:"list"`
		}{CommitList: commitList},
	}
}

// GetBranchList get all branch list
func (Repository) GetBranchList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.GetProjectPath(project.ID)
	git := utils.GIT{Dir: srcPath}

	if err := git.Fetch(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Branch("-r", "--sort=-committerdate"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	unformatBranchList := strings.Split(git.Output.String(), "\n")
	var branchList []string
	for _, row := range unformatBranchList {
		branch := strings.Trim(row, " ")
		if len(branch) != 0 {
			branchList = append(branchList, branch)
		}
	}
	return &core.Response{
		Data: struct {
			BranchList []string `json:"list"`
		}{BranchList: branchList},
	}
}

// GetTagList get latest 10 tag list
func (Repository) GetTagList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.GetProjectPath(project.ID)
	git := utils.GIT{Dir: srcPath}

	if err := git.Add("."); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Reset("--hard"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Pull(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Log("--tags", "-n", "10", "--no-walk", "--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	tagList := utils.ParseGITLog(git.Output.String())

	return &core.Response{
		Data: struct {
			TagList []utils.Commit `json:"list"`
		}{TagList: tagList},
	}
}

