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

	_, err := model.ProjectUser{ProjectID: reqData.ProjectID, UserID: gp.TokenInfo.ID}.GetDataByProjectUser()
	if err != nil {
		return errors.New("无权限进行此操作")
	}
	return nil
}
