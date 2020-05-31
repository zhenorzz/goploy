package ws

import (
	"goploy/model"
)

// ProjectMessage is publish project message struct
type ProjectMessage struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	State       uint8  `json:"state"`
	Message     string `json:"message"`
}

const (
	ProjectFail = 0
	GitClone = 1
	GitReset = 1
	GitSwitchBranch = 2
	GitClean = 3
	GitCheckout = 4
	GitPull = 5
	AfterPullScript = 6
	Rsync = 7
	ProjectSuccess = 8
)




func (projectMessage ProjectMessage) canSendTo(client *Client) error {
	_, err := model.Project{ID: projectMessage.ProjectID}.GetUserProjectData(client.UserInfo.ID, client.UserInfo.Role, client.UserInfo.ManageGroupStr)
	return err
}
