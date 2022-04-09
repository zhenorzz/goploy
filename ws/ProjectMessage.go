package ws

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

func (projectMessage ProjectMessage) canSendTo(*Client) error {
	return nil
}
