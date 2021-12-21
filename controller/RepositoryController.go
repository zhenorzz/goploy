package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/repository"
	"github.com/zhenorzz/goploy/response"
	"net/http"
	"strconv"
)

// Repository struct
type Repository Controller

func (r Repository) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/repository/getCommitList", http.MethodGet, r.GetCommitList),
		core.NewRoute("/repository/getBranchList", http.MethodGet, r.GetBranchList),
		core.NewRoute("/repository/getTagList", http.MethodGet, r.GetTagList),
	}
}

// GetCommitList get latest 10 commit list
func (Repository) GetCommitList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := repo.BranchLog(project.ID, gp.URLQuery.Get("branch"), 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			CommitList []repository.CommitInfo `json:"list"`
		}{CommitList: list},
	}
}

// GetBranchList get all branch list
func (Repository) GetBranchList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	branchList, err := repo.BranchList(project.ID)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			BranchList []string `json:"list"`
		}{BranchList: branchList},
	}
}

// GetTagList get latest 10 tag list
func (Repository) GetTagList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := repo.TagLog(project.ID, 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			TagList []repository.CommitInfo `json:"list"`
		}{TagList: list},
	}
}
