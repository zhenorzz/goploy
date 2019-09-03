package ws

// 单播
import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"goploy/core"
)

// UnicastClient stores a client information
type UnicastClient struct {
	Conn     *websocket.Conn
	UserID   int64
	UserName string
}

// UnicastData is message struct
type UnicastData struct {
	ToUserID int64
	Message  interface{}
}

// ProjectMessage is publish project message struct
type ProjectMessage struct {
	ProjectID  int64  `json:"projectId"`
	ServerID   int64  `json:"serverId"`
	ServerName string `json:"serverName"`
	UserID     int64  `json:"userId"`
	State      uint8  `json:"state"`
	Message    string `json:"message"`
	DataType   uint8  `json:"dataType"`
}

// InstallMessage is server install message struct
type InstallMessage struct {
	ServerID int64  `json:"serverId"`
	UserID   int64  `json:"userId"`
	State    uint8  `json:"state"`
	Message  string `json:"message"`
	DataType uint8  `json:"dataType"`
}

// UnicastHub is a client struct
type UnicastHub struct {
	// Registered clients.
	clients map[*UnicastClient]bool

	// Inbound messages from the clients.
	UnicastData chan *UnicastData

	// Register requests from the clients.
	Register chan *UnicastClient

	// Unregister requests from clients.
	Unregister chan *UnicastClient
}

// ErrorType =>错误信息
// GitType=>git信息
// RsyncType=>rsync信息
// ScriptType => 运行脚本信息
const (
	ErrorType  = 0
	LocalType  = 1
	RsyncType  = 2
	ScriptType = 3
)

var instance *UnicastHub

// GetUnicastHub it will only init once in main.go
func GetUnicastHub() *UnicastHub {
	if instance == nil {
		instance = &UnicastHub{
			UnicastData: make(chan *UnicastData),
			clients:     make(map[*UnicastClient]bool),
			Register:    make(chan *UnicastClient),
			Unregister:  make(chan *UnicastClient),
		}
	}
	return instance
}

// Unicast the publish information in websocket
func (hub *UnicastHub) Unicast(w http.ResponseWriter, gp *core.Goploy) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			}
			return false
		},
	}
	c, err := upgrader.Upgrade(w, gp.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	hub.Register <- &UnicastClient{
		Conn:     c,
		UserID:   gp.TokenInfo.ID,
		UserName: gp.TokenInfo.Name,
	}
}

// Run goroutine run the sync hub
func (hub *UnicastHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Conn.Close()
			}
		case broadcast := <-hub.UnicastData:
			for client := range hub.clients {
				if client.UserID != broadcast.ToUserID {
					continue
				}
				if err := client.Conn.WriteJSON(broadcast.Message); websocket.IsCloseError(err) {
					hub.Unregister <- client
				}
			}
		}
	}
}
