package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Project struct
type Project struct{}

// Add one project
func (project *Project) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Owner      string `json:"owner"`
		Project    string `json:"project"`
		Repository string `json:"repository"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.Project{
		Owner:      reqData.Owner,
		Project:    reqData.Project,
		Repository: reqData.Repository,
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

// Branch list in this project
func (project *Project) Branch(w http.ResponseWriter, r *http.Request) {
	type Branch struct {
		Name string `json:"name"`
	}
	type Branches []Branch
	type RepData struct {
		Branches Branches `json:"branchList"`
	}
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/godis/branches")
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var branches Branches
	err = json.Unmarshal(body, &branches)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Branches: branches}}
	response.Json(w)
}

// Commit list in this Branch
func (project *Project) Commit(w http.ResponseWriter, r *http.Request) {
	type Commits struct {
		Sha    string `json:"sha"`
		NodeID string `json:"node_id"`
		Commit struct {
			Committer struct {
				Name  string `json:"name"`
				Email string `json:"email"`
				Date  string `json:"date"`
			} `json:"committer"`
			Message string `json:"message"`
			Tree    struct {
				Sha string `json:"sha"`
				URL string `json:"url"`
			} `json:"tree"`
		} `json:"commit"`
	}
	type RepData struct {
		Commits []Commits `json:"commitList"`
	}
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/godis/commits?sha=btree")
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var commits []Commits
	err = json.Unmarshal(body, &commits)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Commits: commits}}
	response.Json(w)
}
