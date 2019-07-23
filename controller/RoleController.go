package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Role struct
type Role Controller

// GetOption role list
func (role Role) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Role model.Roles `json:"roleList"`
	}

	roleList, err := model.Role{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Role: roleList}}
	response.JSON(w)
}

// GetPermissionList get role list and permission list
func (role Role) GetPermissionList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		RoleList       model.Roles       `json:"roleList"`
		PermissionTree model.Permissions `json:"permissionTree"`
	}

	roleList, err := model.Role{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	permissions, err := model.Permission{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	var tempPermissions model.Permissions
	for _, permission := range permissions {
		if permission.PID == 0 {
			for _, pmChild := range permissions {
				if pmChild.PID == permission.ID {
					permission.Children = append(permission.Children, pmChild)
				}
			}
			tempPermissions = append(tempPermissions, permission)
		}
	}
	response := core.Response{Data: RepData{RoleList: roleList, PermissionTree: tempPermissions}}
	response.JSON(w)
}

// Edit one Role
func (role Role) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID             uint32 `json:"id"`
		Name           string `json:"name"`
		Remark         string `json:"remark"`
		PermissionList string `json:"permissionList"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Role{
		ID:             reqData.ID,
		Name:           reqData.Name,
		Remark:         reqData.Remark,
		PermissionList: reqData.PermissionList,
		UpdateTime:     time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}
