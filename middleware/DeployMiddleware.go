package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
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

	_, err := model.Project{ID: reqData.ProjectID}.GetUserProjectData(gp.UserInfo.ID)
	if err != nil {
		return errors.New("no permission")
	}
	return nil
}

// HasPublishAuth check the user has publish auth
func FilterEvent(w http.ResponseWriter, gp *core.Goploy) error {

	if XGitHubEvent := gp.Request.Header.Get("X-GitHub-Event"); len(XGitHubEvent) != 0 && XGitHubEvent == "push" {
		return nil
	} else if XGitLabEvent := gp.Request.Header.Get("X-Gitlab-Event"); len(XGitLabEvent) != 0 && XGitLabEvent == "Push Hook" {
		return nil
	} else if XGiteeEvent := gp.Request.Header.Get("X-Gitee-Event");len(XGiteeEvent) != 0 && XGiteeEvent == "Push Hook" {
		return nil
	} else {
		return errors.New("only receive push event")
	}
}
