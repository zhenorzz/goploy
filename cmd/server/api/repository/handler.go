// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repository

import (
	"errors"
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/repo"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"io/fs"
	"net/http"
	"os"
	"path"
)

type Repository api.API

func (r Repository) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/repository/getCommitList", http.MethodGet, r.GetCommitList),
		server.NewRoute("/repository/getBranchList", http.MethodGet, r.GetBranchList),
		server.NewRoute("/repository/getTagList", http.MethodGet, r.GetTagList),
		server.NewRoute("/repository/getFileList", http.MethodGet, r.GetFileList).Permissions(config.ManageRepository),
		server.NewRoute("/repository/previewFile", http.MethodGet, r.PreviewFile).Permissions(config.ManageRepository),
		server.NewRoute("/repository/downloadFile", http.MethodGet, r.DownloadFile).Permissions(config.ManageRepository),
		server.NewRoute("/repository/deleteFile", http.MethodDelete, r.DeleteFile).Permissions(config.ManageRepository),
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

	list, err := r.BranchLog(project.ID, reqData.Branch, 50)
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

// GetFileList lists repository file
// @Summary List repository file
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.GetFileList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.GetFileList.RespData}
// @Router /repository/getFileList [get]
func (Repository) GetFileList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID  int64  `schema:"id" validate:"required,gt=0"`
		Dir string `schema:"dir"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	type item struct {
		Name  string `json:"name"`
		Mode  string `json:"mode"`
		IsDir bool   `json:"isDir"`
	}
	type RespData struct {
		List []item `json:"list"`
	}

	dirEntries, err := os.ReadDir(path.Join(config.GetProjectPath(reqData.ID), reqData.Dir))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return response.JSON{Data: RespData{List: []item{}}}
		}
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	var list []item
	for _, d := range dirEntries {
		i := item{
			Name:  d.Name(),
			IsDir: d.IsDir(),
			Mode:  d.Type().String(),
		}
		list = append(list, i)
	}

	return response.JSON{
		Data: RespData{List: list},
	}
}

// PreviewFile preview repository file
// @Summary Preview repository file
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.PreviewFile.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.PreviewFile.RespData}
// @Router /repository/previewFile [get]
func (Repository) PreviewFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `schema:"id" validate:"required,gt=0"`
		File string `schema:"file" validate:"required,min=1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.File{Filename: path.Join(config.GetProjectPath(reqData.ID), reqData.File), Disposition: "inline"}
}

// DownloadFile download repository file
// @Summary Download repository file
// @Tags Repository
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query repository.DownloadFile.ReqData true "query params"
// @Success 200 {object} response.JSON{data=repository.DownloadFile.RespData}
// @Router /repository/downloadFile [get]
func (Repository) DownloadFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `schema:"id" validate:"required,gt=0"`
		File string `schema:"file" validate:"required,min=1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.File{Filename: path.Join(config.GetProjectPath(reqData.ID), reqData.File), Disposition: "attachment"}
}

// DeleteFile delete repository file
// @Summary Delete repository file
// @Tags Repository
// @Produce json
// @Param request body project.DeleteFile.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /repository/deleteFile [delete]
func (Repository) DeleteFile(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"required,gt=0"`
		File string `json:"file" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	file := path.Join(config.GetProjectPath(reqData.ID), reqData.File)
	fi, err := os.Stat(file)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if fi.IsDir() == true {
		err = os.RemoveAll(file)
	} else {
		err = os.Remove(file)
	}

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
