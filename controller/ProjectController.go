package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

type Project struct {
}

func (project *Project) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Owner      string `json:"owner"`
		Project    string `json:"project"`
		Repository string `json:"repository"`
	}
	type RepData struct {
		Token string `json:"token"`
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

func (project *Project) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Token string `json:"token"`
	}

	model := model.Projects{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: model}
	response.Json(w)
}
