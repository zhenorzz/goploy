package controller

import (
	"net/http"

	"github.com/zhenorzz/goploy/core"
)

// Role struct
type Role Controller

// GetOption role list
func (role Role) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		RoleList []string `json:"roleList"`
	}

	response := core.Response{Data: RepData{RoleList: []string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager, core.RoleMember}}}
	response.JSON(w)
}
