// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

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
	type ReqData struct {
		ID     int64  `schema:"id" validate:"gt=0"`
		Branch string `schema:"branch"  validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: reqData.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := repo.BranchLog(project.ID, reqData.Branch, 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			CommitList []repository.CommitInfo `json:"list"`
		}{CommitList: list},
	}
}

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
