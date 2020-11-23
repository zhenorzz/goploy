package controller

import (
	"database/sql"
	"errors"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/service"
	"github.com/zhenorzz/goploy/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Deploy struct
type Deploy Controller

// GetList -
func (Deploy) GetList(gp *core.Goploy) *core.Response {
	projectName := gp.URLQuery.Get("projectName")
	projects, err := model.Project{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		Name:        projectName,
	}.GetUserProjectList()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Project model.Projects `json:"list"`
		}{Project: projects},
	}
}

// GetPreview deploy detail
func (Deploy) GetPreview(gp *core.Goploy) *core.Response {
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
	return &core.Response{
		Data: struct {
			GitTraceList model.PublishTraces `json:"gitTraceList"`
			Pagination   model.Pagination    `json:"pagination"`
		}{GitTraceList: gitTraceList, Pagination: pagination},
	}
}

// GetDetail deploy detail
func (Deploy) GetDetail(gp *core.Goploy) *core.Response {

	lastPublishToken := gp.URLQuery.Get("lastPublishToken")

	publishTraceList, err := model.PublishTrace{Token: lastPublishToken}.GetListByToken()
	if err == sql.ErrNoRows {
		return &core.Response{Code: core.Error, Message: "No deploy record"}
	} else if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			PublishTraceList model.PublishTraces `json:"publishTraceList"`
		}{PublishTraceList: publishTraceList},
	}
}

// GetCommitList get latest 10 commit list
func (Deploy) GetCommitList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.GetProjectPath(project.ID)
	git := utils.GIT{Dir: srcPath}
	if err := git.Clean("-f"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Checkout("--", "."); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Pull(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Log("--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`", "-n", "10"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	commitList := utils.ParseGITLog(git.Output.String())

	return &core.Response{
		Data: struct {
			CommitList []utils.Commit `json:"commitList"`
		}{CommitList: commitList},
	}
}

// GetTagList get latest 10 tag list
func (Deploy) GetTagList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	srcPath := core.GetProjectPath(project.ID)
	git := utils.GIT{Dir: srcPath}
	if err := git.Clean("-f"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Checkout("--", "."); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Pull(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	if err := git.Log("--tags", "-n", "10", "--no-walk", "--stat", "--pretty=format:`start`%H`%an`%at`%s`%d`"); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error() + ", detail: " + git.Err.String()}
	}

	tagList := utils.ParseGITLog(git.Output.String())

	return &core.Response{
		Data: struct {
			TagList []utils.Commit `json:"tagList"`
		}{TagList: tagList},
	}
}

// Publish the project
func (Deploy) Publish(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64  `json:"projectId" validate:"gt=0"`
		Commit    string `json:"commit"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	var err error
	project, err := model.Project{ID: reqData.ProjectID}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if project.Review == model.Enable && gp.Namespace.Role == core.RoleMember {
		err = projectReview(gp, project, reqData.Commit)
	} else {
		err = projectDeploy(gp, project, reqData.Commit)
	}
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// ResetState -
func (Deploy) ResetState(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64 `json:"projectId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Project{ID: reqData.ProjectID}).ResetState(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}

// GreyPublish the project
func (Deploy) GreyPublish(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectID int64   `json:"projectId" validate:"gt=0"`
		Commit    string  `json:"commit"`
		ServerIDs []int64 `json:"serverIds"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	var err error
	project, err := model.Project{ID: reqData.ProjectID}.GetData()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if project.DeployState == model.ProjectDeploying {
		return &core.Response{Code: core.Error, Message: "project is being build"}
	}

	bindProjectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectServers := model.ProjectServers{}

	for _, projectServer := range bindProjectServers {
		for _, serverID := range reqData.ServerIDs {
			if projectServer.ServerID == serverID {
				projectServers = append(projectServers, projectServer)
			}
		}
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

	return &core.Response{}
}

func (Deploy) Review(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ProjectReviewID int64 `json:"projectReviewId" validate:"gt=0"`
		State           uint8 `json:"state" validate:"gt=0"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectReviewModel := model.ProjectReview{
		ID:       reqData.ProjectReviewID,
		State:    reqData.State,
		Editor:   gp.UserInfo.Name,
		EditorID: gp.UserInfo.ID,
	}
	projectReview, err := projectReviewModel.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if projectReview.State != model.PENDING {
		return &core.Response{Code: core.Error, Message: "Project review state is invalid"}
	}

	if reqData.State == model.APPROVE {
		project, err := model.Project{ID: projectReview.ProjectID}.GetData()
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
		if err := projectDeploy(gp, project, projectReview.CommitID); err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
	}
	projectReviewModel.EditRow()

	return &core.Response{}
}

// Webhook -
func (Deploy) Webhook(gp *core.Goploy) *core.Response {
	projectID, err := strconv.ParseInt(gp.URLQuery.Get("project_id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	// other event is blocked in deployMiddleware
	type ReqData struct {
		Ref string `json:"ref" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: projectID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if project.State != model.Disable {
		return &core.Response{Code: core.Deny, Message: "Project is disabled"}
	}

	if project.AutoDeploy != model.ProjectWebhookDeploy {
		return &core.Response{Code: core.Deny, Message: "Webhook auto deploy turn off, go to project setting turn on"}
	}

	if branch := strings.Split(reqData.Ref, "/")[2]; project.Branch != branch {
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

// Callback -
func (Deploy) Callback(gp *core.Goploy) *core.Response {
	projectReviewID, err := strconv.ParseInt(gp.URLQuery.Get("project_review_id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	projectReviewModel := model.ProjectReview{
		ID:       projectReviewID,
		State:    model.APPROVE,
		Editor:   "admin",
		EditorID: 1,
	}
	projectReview, err := projectReviewModel.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if projectReview.State != model.PENDING {
		return &core.Response{Code: core.Error, Message: "Project review state is invalid"}
	}

	project, err := model.Project{ID: projectReview.ProjectID}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if err := projectDeploy(gp, project, projectReview.CommitID); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	projectReviewModel.EditRow()

	return &core.Response{}
}

func projectDeploy(gp *core.Goploy, project model.Project, commitID string) error {
	if project.DeployState == model.ProjectDeploying {
		return errors.New("project is being build by other")
	}

	projectServers, err := model.ProjectServer{ProjectID: project.ID}.GetBindServerListByProjectID()

	if err != nil {
		return err
	}
	project.PublisherID = gp.UserInfo.ID
	project.PublisherName = gp.UserInfo.Name
	project.DeployState = model.ProjectDeploying
	project.LastPublishToken = uuid.New().String()
	err = project.Publish()
	if err != nil {
		return err
	}
	go service.Sync{
		UserInfo:       gp.UserInfo,
		Project:        project,
		ProjectServers: projectServers,
		CommitID:       commitID,
	}.Exec()
	return nil
}

func projectReview(gp *core.Goploy, project model.Project, commitID string) error {
	if len(commitID) == 0 {
		return errors.New("commit id is required")
	}
	projectReviewModel := model.ProjectReview{
		ProjectID: project.ID,
		CommitID:  commitID,
		Creator:   gp.UserInfo.Name,
		CreatorID: gp.UserInfo.ID,
	}
	reviewURL := project.ReviewURL
	if len(reviewURL) > 0 {
		reviewURL = strings.Replace(reviewURL, "__PROJECT_ID__", strconv.FormatInt(project.ID, 10), 1)
		reviewURL = strings.Replace(reviewURL, "__PROJECT_NAME__", project.Name, 1)
		reviewURL = strings.Replace(reviewURL, "__BRANCH__", project.Branch, 1)
		reviewURL = strings.Replace(reviewURL, "__ENVIRONMENT__", strconv.Itoa(int(project.Environment)), 1)
		reviewURL = strings.Replace(reviewURL, "__COMMIT_ID__", commitID, 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISH_TIME__", strconv.FormatInt(time.Now().Unix(), 10), 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISHER_ID__", gp.UserInfo.Name, 1)
		reviewURL = strings.Replace(reviewURL, "__PUBLISHER_NAME__", strconv.FormatInt(gp.UserInfo.ID, 10), 1)

		projectReviewModel.ReviewURL = reviewURL
	}
	id, err := projectReviewModel.AddRow()
	if err != nil {
		return err
	}
	if len(reviewURL) > 0 {
		callback := "http://"
		if gp.Request.TLS != nil {
			callback = "https://"
		}
		callback += gp.Request.Host + "/deploy/callback?project_review_id=" + strconv.FormatInt(id, 10)
		callback = url.QueryEscape(callback)
		reviewURL = strings.Replace(reviewURL, "__CALLBACK__", callback, 1)

		resp, err := http.Get(reviewURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}
