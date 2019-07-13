package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Project struct
type Project struct{}

// Get project list
func (project *Project) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	model := model.Projects{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Project: model}}
	response.Json(w)
}

// GetDetail project detail
func (project *Project) GetDetail(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		ProjectDetail model.ProjectDetail `json:"projectDetail"`
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.Json(w)
		return
	}
	model := model.ProjectDetail{}
	model.ID = uint32(id)
	if err := model.Detail(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{ProjectDetail: model}}
	response.Json(w)
}

// Add one project
func (project *Project) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Name      string   `json:"name"`
		URL       string   `json:"url"`
		Path      string   `json:"path"`
		ServerIDs []uint32 `json:"serverIds"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectModel := model.Project{
		Name:       reqData.Name,
		URL:        reqData.URL,
		Path:       reqData.Path,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = projectModel.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID:  projectModel.ID,
			ServerID:   serverID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}
	err = projectServersModel.AddMany()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.Json(w)
}

// Create new repository
func (project *Project) Create(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectModel := model.Project{
		ID: reqData.ID,
	}
	err = projectModel.QueryRow()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	projectModel.Status = 1
	err = projectModel.ChangeStatus()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	path := "./repository/" + projectModel.Name
	repo := projectModel.URL

	// clone repository async
	go func(id uint32, path, repo string) {
		projectModel := model.Project{
			ID: id,
		}
		err = os.RemoveAll(path)
		if err != nil {
			projectModel.Status = 3
			_ = projectModel.ChangeStatus()
			fmt.Println(err)
			return
		}
		cmd := exec.Command("git", "clone", repo, path)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			projectModel.Status = 3
			_ = projectModel.ChangeStatus()
			fmt.Println(err)
			return
		}
		projectModel.Status = 2
		_ = projectModel.ChangeStatus()
	}(reqData.ID, path, repo)

	response := core.Response{Message: "初始化中，请稍后"}
	response.Json(w)
}
