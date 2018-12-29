package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Deploy struct
type Deploy struct{}

// Get deploy list
func (deploy *Deploy) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Deploys model.Deploys `json:"deployList"`
	}

	model := model.Deploys{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Deploys: model}}
	response.Json(w)
}

// Publish the project
func (deploy *Deploy) Publish(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Deploys model.Deploys `json:"deployList"`
	}

	model := model.Deploys{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Deploys: model}}
	response.Json(w)
}

// Add one deploy item
func (deploy *Deploy) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ProjectID uint32 `json:"projectId"`
		Branch    string `json:"branch"`
		Commit    string `json:"commit"`
		CommitSha string `json:"commitSha"`
		Type      uint8  `json:"type"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.Deploy{
		ProjectID:  reqData.ProjectID,
		Branch:     reqData.Branch,
		Commit:     reqData.Commit,
		CommitSha:  reqData.CommitSha,
		Type:       reqData.Type,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = model.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.Json(w)
}
