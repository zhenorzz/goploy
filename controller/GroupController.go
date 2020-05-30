package controller

import (
	"encoding/json"
	"goploy/core"
	"goploy/model"
	"net/http"
	"strconv"
	"strings"
)

// Group struct
type Group Controller

// GetList Group list
func (group Group) GetList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Groups     model.Groups     `json:"groupList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	var groupList model.Groups
	if gp.UserInfo.Role == core.RoleAdmin || gp.UserInfo.Role == core.RoleManager {
		groupList, pagination, err = model.Group{}.GetList(pagination)
	} else {
		groupList, pagination, err = model.Group{}.GetListInGroupIDs(strings.Split(gp.UserInfo.ManageGroupStr, ","), pagination)
	}
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{Data: RespData{Groups: groupList, Pagination: pagination}}
}

// GetOption Group list
func (group Group) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Groups model.Groups `json:"groupList"`
	}
	var (
		groupList model.Groups
		err       error
	)
	if gp.UserInfo.Role == core.RoleAdmin || gp.UserInfo.Role == core.RoleManager {
		groupList, err = model.Group{}.GetAll()
	} else {
		groupList, err = model.Group{}.GetAllInGroupIDs(strings.Split(gp.UserInfo.ManageGroupStr, ","))
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Groups: groupList}}
}

// GetDeployOption Group list
func (group Group) GetDeployOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Groups model.Groups `json:"groupList"`
	}
	var (
		groupList model.Groups
		err       error
	)
	if gp.UserInfo.Role == core.RoleAdmin || gp.UserInfo.Role == core.RoleManager {
		groupList, err = model.Group{}.GetAll()
	} else {
		projects, err := model.ProjectUser{
			UserID: gp.UserInfo.ID,
		}.GetListLeftJoinProjectByUserID()
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
		groupIDs := strings.Split(gp.UserInfo.ManageGroupStr, ",")
		for _, project := range projects {
			groupIDs = append(groupIDs, strconv.FormatInt(project.GroupID, 10))
		}
		groupList, err = model.Group{}.GetAllInGroupIDs(groupIDs)
	}

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{Data: RespData{Groups: groupList}}
}

// Add one Group
func (group Group) Add(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	id, err := model.Group{Name: reqData.Name}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit one Group
func (group Group) Edit(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"gt=0"`
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err = model.Group{
		ID:   reqData.ID,
		Name: reqData.Name,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Remove one Server
func (group Group) Remove(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err = model.Group{ID: reqData.ID}.Remove()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}
