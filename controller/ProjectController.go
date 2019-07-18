package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
)

// Project struct
type Project Controller

// GetList project list
func (project Project) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Project model.Projects `json:"projectList"`
	}

	projectList, err := model.Project{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Project: projectList}}
	response.JSON(w)
}

// GetBindServerList project detail
func (project Project) GetBindServerList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		ProjectServers model.ProjectServers `json:"projectServerMap"`
	}
	id, err := strconv.Atoi(gp.URLQuery.Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.JSON(w)
		return
	}
	projectServersMap, err := model.ProjectServer{ProjectID: uint32(id)}.GetBindServerListByProjectID()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{ProjectServers: projectServersMap}}
	response.JSON(w)
}

// GetBindUserList project detail
func (project Project) GetBindUserList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		ProjectUsers model.ProjectUsers `json:"projectUserMap"`
	}
	id, err := strconv.Atoi(gp.URLQuery.Get("id"))
	if err != nil {
		response := core.Response{Code: 1, Message: "id参数错误"}
		response.JSON(w)
		return
	}
	projectUsersMap, err := model.ProjectUser{ProjectID: uint32(id)}.GetBindUserListByProjectID()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{ProjectUsers: projectUsersMap}}
	response.JSON(w)
}

// Add one project
func (project Project) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		Path        string   `json:"path"`
		Script      string   `json:"script"`
		RsyncOption string   `json:"rsyncOption"`
		ServerIDs   []uint32 `json:"serverIds"`
		UserIDs     []uint32 `json:"userIds"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	if _, err := utils.ParseCommandLine(reqData.RsyncOption); err != nil {
		response := core.Response{Code: 1, Message: "Rsync Option错误，请输入正确的参数格式"}
		response.JSON(w)
		return
	}

	projectID, err := model.Project{
		Name:        reqData.Name,
		URL:         reqData.URL,
		Path:        reqData.Path,
		Script:      reqData.Script,
		RsyncOption: reqData.RsyncOption,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID:  projectID,
			ServerID:   serverID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}

	if err := projectServersModel.AddMany(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID:  projectID,
			UserID:     userID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	if err := projectUsersModel.AddMany(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one Project
func (project Project) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID          uint32 `json:"id"`
		Name        string `json:"name"`
		URL         string `json:"url"`
		Path        string `json:"path"`
		Script      string `json:"script"`
		RsyncOption string `json:"rsyncOption"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Project{
		ID:          reqData.ID,
		Name:        reqData.Name,
		URL:         reqData.URL,
		Path:        reqData.Path,
		Script:      reqData.Script,
		RsyncOption: reqData.RsyncOption,
		UpdateTime:  time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}

// AddServer one project
func (project Project) AddServer(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectID uint32   `json:"projectId"`
		ServerIDs []uint32 `json:"serverIds"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	projectID := reqData.ProjectID

	projectServersModel := model.ProjectServers{}
	for _, serverID := range reqData.ServerIDs {
		projectServerModel := model.ProjectServer{
			ProjectID:  projectID,
			ServerID:   serverID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectServersModel = append(projectServersModel, projectServerModel)
	}

	if err := projectServersModel.AddMany(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// AddUser one project
func (project Project) AddUser(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectID uint32   `json:"projectId"`
		UserIDs   []uint32 `json:"userIds"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	projectID := reqData.ProjectID

	projectUsersModel := model.ProjectUsers{}
	for _, userID := range reqData.UserIDs {
		projectUserModel := model.ProjectUser{
			ProjectID:  projectID,
			UserID:     userID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}

	if err := projectUsersModel.AddMany(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// RemoveProjectServer one Project
func (project Project) RemoveProjectServer(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectServerID uint32 `json:"projectServerId"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.ProjectServer{
		ID: reqData.ProjectServerID,
	}.DeleteRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}

// RemoveProjectUser one Project
func (project Project) RemoveProjectUser(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ProjectUserID uint32 `json:"projectUserId"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.ProjectUser{
		ID: reqData.ProjectUserID,
	}.DeleteRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}
