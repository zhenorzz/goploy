package ws

// ProjectMessage is publish project message struct
type ServerTemplateMessage struct {
	ServerID     int64  `json:"serverId"`
	ServerName   string `json:"serverName"`
	Detail       string `json:"detail"`
	Ext          string `json:"ext"`
}

func (serverTemplateMessage ServerTemplateMessage) canSendTo(client *Client) error {
	return nil
}
