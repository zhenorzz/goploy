package controller

import (
	"net/http"

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
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Role: roleList}}
	response.Json(w)
}
