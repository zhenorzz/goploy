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
	TaskWaiting     = 1
	RepoFollow      = 2
	AfterPullScript = 3
	Rsync           = 4
	DeploySuccess   = 5
)

func (projectMessage ProjectMessage) canSendTo(client *Client) error {
	_, err := model.Project{ID: projectMessage.ProjectID, UserID: client.UserInfo.ID}.GetUserProjectData()
	return err
}
