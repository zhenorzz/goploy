package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"goploy/core"
	"goploy/model"
)

// Group struct
type Group Controller

// GetList Group list
func (group Group) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Groups     model.Groups     `json:"groupList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	var (
		groupList model.Groups
	)
	if gp.UserInfo.Role == core.RoleAdmin || gp.UserInfo.Role == core.RoleManager {
		groupList, pagination, err = model.Group{}.GetList(pagination)
	} else {
		groupList, pagination, err = model.Group{}.GetListInGroupIDs(strings.Split(gp.UserInfo.ManageGroupStr, ","), pagination)
	}
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{Groups: groupList, Pagination: pagination}}
	response.JSON(w)
}

// GetOption Group list
func (group Group) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Groups model.Groups `json:"groupList"`
	}

	groupList, err := model.Group{}.GetAll()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{Groups: groupList}}
	response.JSON(w)
}

// Add one Group
func (group Group) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name string `json:"name"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err = model.Group{
		Name:       reqData.Name,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one Group
func (group Group) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Group{
		ID:         reqData.ID,
		Name:       reqData.Name,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}

// Remove one Server
func (group Group) Remove(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Group{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}
