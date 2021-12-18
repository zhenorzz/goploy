package ws

type MonitorMessage struct {
	MonitorID    int64  `json:"monitorId"`
	State        uint8  `json:"state"`
	ErrorContent string `json:"errorContent"`
}

const (
	MonitorTurnOff = 0
)

func (projectMessage MonitorMessage) canSendTo(client *Client) error {
	return nil
}
