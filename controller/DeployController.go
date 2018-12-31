package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &reqData); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	deployModel := model.Deploy{
		ID: reqData.ID,
	}

	if err := deployModel.QueryRow(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	projectModel := model.Project{
		ID: deployModel.ProjectID,
	}

	if err := projectModel.QueryRow(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if projectModel.Status != 2 {
		response := core.Response{Code: 1, Message: "项目尚未初始化"}
		response.Json(w)
		return
	}
	srcPath := "./repository/" + projectModel.Owner + "/" + projectModel.Repository
	descPath := "root@129.204.80.253:/home/ubuntu/data"
	cmd := exec.Command("rsync", "-rtv", srcPath, descPath)
	// cmd.Stderr = os.Stderr
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, errbuf.String())
		return
	}
	response := core.Response{Message: "部署中，请稍后"}
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
