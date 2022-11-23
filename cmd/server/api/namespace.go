// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/model"
	"net/http"
	"strconv"
)

type Namespace API

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

func (Namespace) GetList(gp *server.Goploy) server.Response {
	namespaceList, err := model.Namespace{UserID: gp.UserInfo.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			Namespaces model.Namespaces `json:"list"`
		}{Namespaces: namespaceList},
	}
}

func (Namespace) GetOption(gp *server.Goploy) server.Response {
	namespaceUsers, err := model.NamespaceUser{UserID: gp.UserInfo.ID}.GetUserNamespaceList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) GetUserOption(gp *server.Goploy) server.Response {
	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllUserByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) GetBindUserList(gp *server.Goploy) server.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	namespaceUsers, err := model.NamespaceUser{NamespaceID: id}.GetBindUserListByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := model.Namespace{Name: reqData.Name}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{NamespaceID: id}).AddAdminByNamespaceID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Namespace) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"gt=0"`
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Namespace) AddUser(gp *server.Goploy) server.Response {
	type ReqData struct {
		NamespaceID int64   `json:"namespaceId" validate:"gt=0"`
		UserIDs     []int64 `json:"userIds" validate:"required"`
		RoleID      int64   `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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

func (Namespace) RemoveUser(gp *server.Goploy) server.Response {
	type ReqData struct {
		NamespaceUserID int64 `json:"namespaceUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{ID: reqData.NamespaceUserID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
