package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/response"
	"net/http"
)

type Role Controller

func (r Role) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/role/getList", http.MethodGet, r.GetList).Permissions(permission.ShowRolePage),
		core.NewRoute("/role/getOption", http.MethodGet, r.GetOption),
		core.NewRoute("/role/getPermissionList", http.MethodGet, r.GetPermissionList).Permissions(permission.ShowRolePage),
		core.NewRoute("/role/getPermissionBindings", http.MethodGet, r.GetPermissionBindings).Permissions(permission.ShowRolePage),
		core.NewRoute("/role/add", http.MethodPost, r.Add).Permissions(permission.AddRole),
		core.NewRoute("/role/edit", http.MethodPut, r.Edit).Permissions(permission.EditRole),
		core.NewRoute("/role/remove", http.MethodDelete, r.Remove).Permissions(permission.DeleteRole),
		core.NewRoute("/role/changePermission", http.MethodPut, r.ChangePermission).Permissions(permission.EditPermission),
	}
}

func (Role) GetList(*core.Goploy) core.Response {
	roleList, err := model.Role{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model.Roles `json:"list"`
		}{List: roleList},
	}
}

func (Role) GetOption(*core.Goploy) core.Response {
	list, err := model.Role{}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Data: struct {
		List model.Roles `json:"list"`
	}{list}}
}

func (Role) GetPermissionList(*core.Goploy) core.Response {
	list, err := model.Permission{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model.Permissions `json:"list"`
		}{List: list},
	}
}

func (Role) GetPermissionBindings(gp *core.Goploy) core.Response {
	type ReqData struct {
		RoleID int64 `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.RolePermission{RoleID: reqData.RoleID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			List model.RolePermissions `json:"list"`
		}{List: list},
	}
}

func (Role) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Role{Name: reqData.Name, Description: reqData.Description}.AddRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Role) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.Role{ID: reqData.ID, Name: reqData.Name, Description: reqData.Description}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Role) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.Role{ID: reqData.ID}.DeleteRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Role) ChangePermission(gp *core.Goploy) core.Response {
	type ReqData struct {
		PermissionIDs []int64 `json:"permissionIds" validate:"required"`
		RoleID        int64   `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
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
