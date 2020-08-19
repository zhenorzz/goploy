package controller

import (
	"database/sql"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/utils"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Deploy struct
type Deploy Controller

// GetList -
func (deploy Deploy) GetList(gp *core.Goploy) *core.Response {
	type RespData struct {
		Project model.Projects `json:"list"`
	}
	projectName := gp.URLQuery.Get("projectName")
	projects, err := model.Project{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		Name:        projectName,
	}.GetUserProjectList()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Project: projects}}
}

// GetPreview deploy detail
func (deploy Deploy) GetPreview(gp *core.Goploy) *core.Response {
	type RespData struct {
		GitTraceList model.PublishTraces `json:"gitTraceList"`
		Pagination   model.Pagination    `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectID, err := strconv.ParseInt(gp.URLQuery.Get("projectId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	userID, err := strconv.ParseInt(gp.URLQuery.Get("userId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	state, err := strconv.ParseInt(gp.URLQuery.Get("state"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	gitTraceList, pagination, err := model.PublishTrace{
		ProjectID:    projectID,
		PublisherID:  userID,
		PublishState: int(state),
	}.GetPreview(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{GitTraceList: gitTraceList, Pagination: pagination}}
}

// GetDetail deploy detail
func (deploy Deploy) GetDetail(gp *core.Goploy) *core.Response {
	type RespData struct {
		PublishTraceList model.PublishTraces `json:"publishTraceList"`
	}

	lastPublishToken := gp.URLQuery.Get("lastPublishToken")

	publishTraceList, err := model.PublishTrace{Token: lastPublishToken}.GetListByToken()
	if err == sql.ErrNoRows {
		return &core.Response{Code: core.Error, Message: "No deploy record"}
	} else if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{PublishTraceList: publishTraceList}}
}

// GetCommitList get latest 10 commit list
func (deploy Deploy) GetCommitList(gp *core.Goploy) *core.Response {
	type RespData struct {
		CommitList []utils.Commit `json:"commitList"`
	}

	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.RepositoryPath + project.Name
	git := utils.GIT{Dir: srcPath}
	if err := git.Clean([]string{"-f"}); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Checkout([]string{"--", "."}); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Pull([]string{}); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Log([]string{"--stat", "--pretty=format:`start`%H`%an`%at`%s`", "-n", "10"}); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	commitList := utils.ParseGITLog(git.Output.String())

	return &core.Response{Data: RespData{CommitList: commitList}}
}

// Publish the project
func (deploy Deploy) Publish(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Commit    string `json:"commit"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{
		ID: reqData.ProjectID,
	}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if project.DeployState == model.ProjectDeploying {
		return &core.Response{Code: core.Deny, Message: "Project is being build by other"}
	}

	projectServers, err := model.ProjectServer{ProjectID: reqData.ProjectID}.GetBindServerListByProjectID()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	go service.Sync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
		CommitID:       reqData.Commit,
	}.Exec()
	return &core.Response{Message: "deploying"}
}

// Webhook connect
func (deploy Deploy) Webhook(gp *core.Goploy) *core.Response {
	projectName := gp.URLQuery.Get("project_name")
	// other event is blocked in deployMiddleware
	type ReqData struct {
		Ref string `json:"ref" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	branch := strings.Split(reqData.Ref, "/")[2]

	project, err := model.Project{
		Name: projectName,
	}.GetDataByName()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if project.State != model.Disable {
		return &core.Response{Code: core.Deny, Message: "Project is disabled"}
	}

	if project.AutoDeploy != model.ProjectWebhookDeploy {
		return &core.Response{Code: core.Deny, Message: "Webhook auto deploy turn off, go to project setting turn on"}
	}

	if project.Branch != branch {
		return &core.Response{Code: core.Deny, Message: "Receive branch:" + branch + " push event, not equal to current branch"}
	}

	if project.DeployState == model.ProjectDeploying {
		return &core.Response{Code: core.Deny, Message: "Project is being build by other"}
	}

	gp.UserInfo, err = model.User{ID: 1}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	go service.Sync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
		CommitID:       "",
	}.Exec()
	return &core.Response{Message: "receive push signal"}
}
