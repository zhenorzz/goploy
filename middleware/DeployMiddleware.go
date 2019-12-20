package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"goploy/core"
	"goploy/model"
)

// HasPublishAuth check the user has publish auth
func HasPublishAuth(w http.ResponseWriter, gp *core.Goploy) error {
	type ReqData struct {
		ProjectID int64 `json:"projectId"`
	}
	var reqData ReqData
	if err := json.Unmarshal(gp.Body, &reqData); err != nil {
		return err
	}

	_, err := model.Project{ID: reqData.ProjectID}.GetUserProjectData(gp.UserInfo.ID, gp.UserInfo.Role, gp.UserInfo.ManageGroupStr)
	if err != nil {
		return errors.New("no permission")
	}
	return nil
}

// HasPublishAuth check the user has publish auth
func FilterEvent(w http.ResponseWriter, gp *core.Goploy) error {
	XGitHubEvent := gp.Request.Header.Get("X-GitHub-Event")
	if len(XGitHubEvent) != 0 && XGitHubEvent == "push" {
		return nil
	}

	XGitLabEvent := gp.Request.Header.Get("X-Gitlab-Event")
	if len(XGitLabEvent) != 0 && XGitLabEvent == "Push Hook" {
		return nil
	}

	XGiteeEvent := gp.Request.Header.Get("X-Gitee-Event")
	if len(XGiteeEvent) != 0 && XGiteeEvent == "Push Hook" {
		return nil
	}

	return errors.New("")
}
