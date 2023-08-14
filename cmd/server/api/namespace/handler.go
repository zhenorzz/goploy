// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package namespace

import (
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Namespace api.API

func (n Namespace) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/namespace/getList", http.MethodGet, n.GetList).Permissions(config.ShowNamespacePage),
		server.NewRoute("/namespace/getOption", http.MethodGet, n.GetOption),
		server.NewRoute("/namespace/getBindUserList", http.MethodGet, n.GetBindUserList).Permissions(config.ShowNamespacePage),
		server.NewRoute("/namespace/getUserOption", http.MethodGet, n.GetUserOption),
		server.NewRoute("/namespace/add", http.MethodPost, n.Add).Permissions(config.AddNamespace).LogFunc(middleware.AddOPLog),
		server.NewRoute("/namespace/edit", http.MethodPut, n.Edit).Permissions(config.EditNamespace).LogFunc(middleware.AddOPLog),
		server.NewRoute("/namespace/addUser", http.MethodPost, n.AddUser).Permissions(config.AddNamespaceUser).LogFunc(middleware.AddOPLog),
		server.NewRoute("/namespace/removeUser", http.MethodDelete, n.RemoveUser).Permissions(config.DeleteNamespaceUser).LogFunc(middleware.AddOPLog),
	}
}

// GetList lists namespaces
// @Summary List namespaces
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=namespace.GetList.RespData}
// @Router /namespace/getList [get]
func (Namespace) GetList(gp *server.Goploy) server.Response {
	namespaceList, err := model.Namespace{UserID: gp.UserInfo.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		Namespaces model.Namespaces `json:"list"`
	}
	return response.JSON{
		Data: RespData{Namespaces: namespaceList},
	}
}

// GetOption lists namespaces by user id
// @Summary List namespaces by user id
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=namespace.GetOption.RespData}
// @Router /namespace/getOption [get]
func (Namespace) GetOption(gp *server.Goploy) server.Response {
	namespaceUsers, err := model.NamespaceUser{UserID: gp.UserInfo.ID}.GetUserNamespaceList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		NamespaceUsers model.NamespaceUsers `json:"list"`
	}
	return response.JSON{
		Data: RespData{NamespaceUsers: namespaceUsers},
	}
}

// GetUserOption lists namespaces by namespace id
// @Summary List namespaces by namespace id
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=namespace.GetUserOption.RespData}
// @Router /namespace/getUserOption [get]
func (Namespace) GetUserOption(gp *server.Goploy) server.Response {
	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllUserByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		NamespaceUsers model.NamespaceUsers `json:"list"`
	}
	return response.JSON{
		Data: RespData{NamespaceUsers: namespaceUsers},
	}
}

// GetBindUserList lists namespaces by namespace id
// @Summary List namespaces by namespace id
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request query namespace.GetBindUserList.ReqData true "query params"
// @Success 200 {object} response.JSON{data=namespace.GetBindUserList.RespData}
// @Router /namespace/getBindUserList [get]
func (Namespace) GetBindUserList(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	namespaceUsers, err := model.NamespaceUser{NamespaceID: reqData.ID}.GetBindUserListByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		NamespaceUsers model.NamespaceUsers `json:"list"`
	}
	return response.JSON{
		Data: RespData{NamespaceUsers: namespaceUsers},
	}
}

// Add adds a namespace
// @Summary Add a namespace
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body namespace.Add.ReqData true "body params"
// @Success 200 {object} response.JSON{data=namespace.Add.RespData}
// @Router /namespace/add [post]
func (Namespace) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := model.Namespace{Name: reqData.Name}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{NamespaceID: id}).AddAdminByNamespaceID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// Edit edits the namespace
// @Summary Edit the namespace
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body namespace.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /namespace/edit [put]
func (Namespace) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"required,gt=0"`
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Namespace{
		ID:   reqData.ID,
		Name: reqData.Name,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

// AddUser adds a user
// @Summary Add a user
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body namespace.AddUser.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /namespace/addUser [post]
func (Namespace) AddUser(gp *server.Goploy) server.Response {
	type ReqData struct {
		NamespaceID int64   `json:"namespaceId" validate:"required,gt=0"`
		UserIDs     []int64 `json:"userIds" validate:"required"`
		RoleID      int64   `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	namespaceUsersModel := model.NamespaceUsers{}
	for _, userID := range reqData.UserIDs {
		namespaceUserModel := model.NamespaceUser{
			NamespaceID: reqData.NamespaceID,
			UserID:      userID,
			RoleID:      reqData.RoleID,
		}
		namespaceUsersModel = append(namespaceUsersModel, namespaceUserModel)
	}

	if err := namespaceUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

// RemoveUser removes the user
// @Summary Remove the user
// @Tags Namespace
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body namespace.RemoveUser.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /namespace/removeUser [delete]
func (Namespace) RemoveUser(gp *server.Goploy) server.Response {
	type ReqData struct {
		NamespaceUserID int64 `json:"namespaceUserId" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{ID: reqData.NamespaceUserID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
