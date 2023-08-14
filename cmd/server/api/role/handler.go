// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package role

import (
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Role api.API

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

// GetList lists roles
// @Summary List roles
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=role.GetList.RespData}
// @Router /role/getList [get]
func (Role) GetList(*server.Goploy) server.Response {
	roleList, err := model.Role{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		List model.Roles `json:"list"`
	}

	return response.JSON{
		Data: RespData{List: roleList},
	}
}

// GetOption lists all roles
// @Summary List all roles
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=role.GetOption.RespData}
// @Router /role/getOption [get]
func (Role) GetOption(*server.Goploy) server.Response {
	list, err := model.Role{}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		List model.Roles `json:"list"`
	}
	return response.JSON{Data: RespData{list}}
}

// GetPermissionList lists permissions
// @Summary List permissions
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=role.GetPermissionList.RespData}
// @Router /role/getPermissionList [get]
func (Role) GetPermissionList(*server.Goploy) server.Response {
	list, err := model.Permission{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	type RespData struct {
		List model.Permissions `json:"list"`
	}
	return response.JSON{
		Data: RespData{List: list},
	}
}

// GetPermissionBindings lists role permissions
// @Summary List role permissions
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=role.GetPermissionBindings.RespData}
// @Router /role/getPermissionBindings [get]
func (Role) GetPermissionBindings(gp *server.Goploy) server.Response {
	type ReqData struct {
		RoleID int64 `json:"roleId" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.RolePermission{RoleID: reqData.RoleID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		List model.RolePermissions `json:"list"`
	}
	return response.JSON{
		Data: RespData{List: list},
	}
}

// Add adds a role
// @Summary Add a role
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body role.Add.ReqData true "body params"
// @Success 200 {object} response.JSON{data=role.Add.RespData}
// @Router /role/add [post]
func (Role) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Role{Name: reqData.Name, Description: reqData.Description}.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// Edit edits the role
// @Summary Edit the role
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body role.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /role/edit [put]
func (Role) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"required,gt=0"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.Role{ID: reqData.ID, Name: reqData.Name, Description: reqData.Description}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

// Remove removes the role
// @Summary Remove the role
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body role.Remove.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /role/remove [delete]
func (Role) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	namespaceUser, err := (model.NamespaceUser{RoleID: reqData.ID}).GetDataByRoleID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if namespaceUser.ID > 0 {
		return response.JSON{Code: response.Error, Message: "The role has binding user"}
	}

	if err := (model.Role{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

// ChangePermission changes the role permissions
// @Summary Change the role permissions
// @Tags Role
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body role.ChangePermission.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /role/changePermission [put]
func (Role) ChangePermission(gp *server.Goploy) server.Response {
	type ReqData struct {
		PermissionIDs []int64 `json:"permissionIds" validate:"required"`
		RoleID        int64   `json:"roleId" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.RolePermission{RoleID: reqData.RoleID}).DeleteByRoleID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	rolePermissionsModel := model.RolePermissions{}
	for _, PermissionID := range reqData.PermissionIDs {
		rolePermissionModel := model.RolePermission{
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
