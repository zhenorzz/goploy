// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Repository api.API

func (r Repository) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/repository/getCommitList", http.MethodGet, r.GetCommitList),
		server.NewRoute("/repository/getBranchList", http.MethodGet, r.GetBranchList),
		server.NewRoute("/repository/getTagList", http.MethodGet, r.GetTagList),
	}
}

// GetBranchList lists branches
// @Summary List branches
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.GetBranchList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.GetBranchList.RespData}
// @Router /repository/getBranchList [get]
func (Repository) GetBranchList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `schema:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	branchList, err := r.BranchList(project.ID)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		BranchList []string `json:"list"`
	}
	return response.JSON{
		Data: RespData{BranchList: branchList},
	}
}

// GetCommitList lists latest 10 commits
// @Summary List latest 10 commits
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.GetCommitList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.GetCommitList.RespData}
// @Router /repository/getCommitList [get]
func (Repository) GetCommitList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID     int64  `schema:"id" validate:"required,gt=0"`
		Branch string `schema:"branch"  validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	type RespData struct {
		CommitList []repo.CommitInfo `json:"list"`
	}
	return response.JSON{
		Data: RespData{CommitList: list},
	}
}

// GetTagList lists latest 10 tags
// @Summary List latest 10 tags
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.GetTagList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.GetTagList.RespData}
// @Router /repository/getTagList [get]
func (Repository) GetTagList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `schema:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

	list, err := r.TagLog(project.ID, 10)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		TagList []repo.CommitInfo `json:"list"`
	}
	return response.JSON{
		Data: RespData{TagList: list},
	}
}
