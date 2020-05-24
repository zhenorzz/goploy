package ws

import (
	"goploy/model"
)

// ProjectMessage is publish project message struct
type ProjectMessage struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	State       uint8  `json:"state"`
	Message     string `json:"message"`
}

func (projectMessage ProjectMessage) canSendTo(client *Client) error {
	_, err := model.Project{ID: projectMessage.ProjectID}.GetUserProjectData(client.UserInfo.ID, client.UserInfo.Role, client.UserInfo.ManageGroupStr)
	return err
}
