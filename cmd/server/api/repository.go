// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/model"
	"net/http"
	"strconv"
)

// Repository struct
type Repository API

func (r Repository) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/repository/getCommitList", http.MethodGet, r.GetCommitList),
		server.NewRoute("/repository/getBranchList", http.MethodGet, r.GetBranchList),
		server.NewRoute("/repository/getTagList", http.MethodGet, r.GetTagList),
	}
}

// GetCommitList get latest 10 commit list
func (Repository) GetCommitList(gp *server.Goploy) server.Response {
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

	r, err := repo.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := r.BranchLog(project.ID, reqData.Branch, 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			CommitList []repo.CommitInfo `json:"list"`
		}{CommitList: list},
	}
}

func (Repository) GetBranchList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	r, err := repo.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	branchList, err := r.BranchList(project.ID)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			BranchList []string `json:"list"`
		}{BranchList: branchList},
	}
}

func (Repository) GetTagList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	r, err := repo.GetRepo(project.RepoType)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := r.TagLog(project.ID, 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			TagList []repo.CommitInfo `json:"list"`
		}{TagList: list},
	}
}
