// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	model2 "github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Role API

func (r Role) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/role/getList", http.MethodGet, r.GetList).Permissions(config.ShowRolePage),
		server.NewRoute("/role/getOption", http.MethodGet, r.GetOption),
		server.NewRoute("/role/getPermissionList", http.MethodGet, r.GetPermissionList).Permissions(config.ShowRolePage),
		server.NewRoute("/role/getPermissionBindings", http.MethodGet, r.GetPermissionBindings).Permissions(config.ShowRolePage),
		server.NewRoute("/role/add", http.MethodPost, r.Add).Permissions(config.AddRole).LogFunc(middleware.AddOPLog),
		server.NewRoute("/role/edit", http.MethodPut, r.Edit).Permissions(config.EditRole).LogFunc(middleware.AddOPLog),
		server.NewRoute("/role/remove", http.MethodDelete, r.Remove).Permissions(config.DeleteRole).LogFunc(middleware.AddOPLog),
		server.NewRoute("/role/changePermission", http.MethodPut, r.ChangePermission).Permissions(config.EditPermission).LogFunc(middleware.AddOPLog),
	}
}

func (Role) GetList(*server.Goploy) server.Response {
	roleList, err := model2.Role{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model2.Roles `json:"list"`
		}{List: roleList},
	}
}

func (Role) GetOption(*server.Goploy) server.Response {
	list, err := model2.Role{}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Data: struct {
		List model2.Roles `json:"list"`
	}{list}}
}

func (Role) GetPermissionList(*server.Goploy) server.Response {
	list, err := model2.Permission{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model2.Permissions `json:"list"`
		}{List: list},
	}
}

func (Role) GetPermissionBindings(gp *server.Goploy) server.Response {
	type ReqData struct {
		RoleID int64 `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model2.RolePermission{RoleID: reqData.RoleID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model2.RolePermissions `json:"list"`
		}{List: list},
	}
}

func (Role) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model2.Role{Name: reqData.Name, Description: reqData.Description}.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Role) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model2.Role{ID: reqData.ID, Name: reqData.Name, Description: reqData.Description}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Role) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	namespaceUser, err := (model2.NamespaceUser{RoleID: reqData.ID}).GetDataByRoleID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if namespaceUser.ID > 0 {
		return response.JSON{Code: response.Error, Message: "The role has binding user"}
	}

	if err := (model2.Role{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Role) ChangePermission(gp *server.Goploy) server.Response {
	type ReqData struct {
		PermissionIDs []int64 `json:"permissionIds" validate:"required"`
		RoleID        int64   `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model2.RolePermission{RoleID: reqData.RoleID}).DeleteByRoleID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	rolePermissionsModel := model2.RolePermissions{}
	for _, PermissionID := range reqData.PermissionIDs {
		rolePermissionModel := model2.RolePermission{
			PermissionID: PermissionID,
			RoleID:       reqData.RoleID,
		}
		rolePermissionsModel = append(rolePermissionsModel, rolePermissionModel)
	}

	if err := rolePermissionsModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
