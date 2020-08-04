package ws

// ProjectMessage is publish project message struct
type ServerTemplateMessage struct {
	ServerID   int64       `json:"serverId"`
	ServerName string      `json:"serverName"`
	Detail     string      `json:"detail"`
	Ext        interface{} `json:"ext"`
	Type       uint8       `json:"type"`
}

const (
	ServerTemplateRsync  = 1
	ServerTemplateSSH    = 2
	ServerTemplateScript = 3
)

func (serverTemplateMessage ServerTemplateMessage) canSendTo(client *Client) error {
	return nil
}
