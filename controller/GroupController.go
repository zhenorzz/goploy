package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Group struct
type Group Controller

// GetList Group list
func (group Group) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Groups model.Groups `json:"groupList"`
	}

	groupList, err := model.Group{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Groups: groupList}}
	response.JSON(w)
}

// GetOption Group list
func (group Group) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Groups model.Groups `json:"groupList"`
	}

	groupList, err := model.Group{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Groups: groupList}}
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
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err = model.Group{
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

// Edit one Group
func (group Group) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
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
		response := core.Response{Code: 1, Message: err.Error()}
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
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Group{
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
