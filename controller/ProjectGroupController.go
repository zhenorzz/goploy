package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// ProjectGroup struct
type ProjectGroup Controller

// GetList ProjectGroup list
func (projectGroup ProjectGroup) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		ProjectGroups model.ProjectGroups `json:"projectGroupList"`
	}

	projectGroupList, err := model.ProjectGroup{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{ProjectGroups: projectGroupList}}
	response.JSON(w)
}

// GetOption ProjectGroup list
func (projectGroup ProjectGroup) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		ProjectGroups model.ProjectGroups `json:"projectGroupList"`
	}

	projectGroupList, err := model.ProjectGroup{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{ProjectGroups: projectGroupList}}
	response.JSON(w)
}

// Add one ProjectGroup
func (projectGroup ProjectGroup) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name string `json:"name"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err = model.ProjectGroup{
		Name:       reqData.Name,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one ProjectGroup
func (projectGroup ProjectGroup) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID   uint32 `json:"id"`
		Name string `json:"name"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.ProjectGroup{
		ID:         reqData.ID,
		Name:       reqData.Name,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}

// Remove one Server
func (projectGroup ProjectGroup) Remove(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.ProjectGroup{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}
