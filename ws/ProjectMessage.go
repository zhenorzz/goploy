package ws

import (
	"github.com/zhenorzz/goploy/model"
)

// ProjectMessage is publish project message struct
type ProjectMessage struct {
	ProjectID   int64       `json:"projectId"`
	ProjectName string      `json:"projectName"`
	State       uint8       `json:"state"`
	Message     string      `json:"message"`
	Ext         interface{} `json:"ext"`
}

const (
	DeployFail      = 0
	GitCreate       = 1
	GitClean        = 2
	GitFetch        = 3
	GitCheckout     = 4
	GitDone         = 5
	AfterPullScript = 6
	Rsync           = 7
	DeploySuccess   = 8
)

func (projectMessage ProjectMessage) canSendTo(client *Client) error {
	_, err := model.Project{ID: projectMessage.ProjectID, UserID: client.UserInfo.ID}.GetUserProjectData()
	return err
}
