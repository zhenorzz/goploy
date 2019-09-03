package controller

import (
	"net/http"

	"goploy/core"
)

// Role struct
type Role Controller

// GetOption role list
func (role Role) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		RoleList []string `json:"roleList"`
	}

	response := core.Response{Data: RespData{RoleList: []string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager, core.RoleMember}}}
	response.JSON(w)
}
