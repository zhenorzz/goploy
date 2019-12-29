package controller

import (
	"net/http"

	"goploy/core"
)

// Role struct
type Role Controller

// GetOption role list
func (role Role) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		RoleList []string `json:"roleList"`
	}
	return &core.Response{Data: RespData{RoleList: []string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager, core.RoleMember}}}
}
