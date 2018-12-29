package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

// Branch list in this project
func (project *Project) Branch(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID uint32 `json:"id"`
	}
	type Branch struct {
		Name string `json:"name"`
	}
	type Branches []Branch
	type RepData struct {
		Branches Branches `json:"branchList"`
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
	url := "https://api.github.com/repos/" + projectModel.Owner + "/" + projectModel.Project + "/branches"
	resp, err := http.Get(url)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
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
	type ReqData struct {
		ID     uint32 `json:"id"`
		Branch string `json:"branch"`
	}
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
	url := "https://api.github.com/repos/" + projectModel.Owner + "/" + projectModel.Project + "/commits?sha=" + reqData.Branch
	resp, err := http.Get(url)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
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

	path := "./repository/" + projectModel.Owner + "/" + projectModel.Repository
	repo := "https://github.com/" + projectModel.Owner + "/" + projectModel.Repository + ".git"

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
