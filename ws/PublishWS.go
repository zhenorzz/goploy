package ws

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// SyncClient stores a client information
type SyncClient struct {
	Conn       *websocket.Conn
	UserID     uint32
	UserName   string
	ProjectMap map[uint32]struct{}
}

// SyncBroadcast is message struct
type SyncBroadcast struct {
	ProjectID  uint32 `json:"projectId"`
	ServerID   uint32 `json:"serverId"`
	ServerName string `json:"serverName"`
	UserID     uint32 `json:"userId"`
	State      uint8  `json:"state"`
	Message    string `json:"message"`
	DataType   uint8  `json:"dataType"`
}

// SyncHub is a client struct
type SyncHub struct {
	// Registered clients.
	clients map[*SyncClient]bool

	// Inbound messages from the clients.
	Broadcast chan *SyncBroadcast

	// Register requests from the clients.
	Register chan *SyncClient

	// Unregister requests from clients.
	Unregister chan *SyncClient
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

// State
const (
	Fail    = 0
	Success = 1
)

var instance *SyncHub

// GetSyncHub it will only init once in main.go
func GetSyncHub() *SyncHub {
	if instance == nil {
		instance = &SyncHub{
			Broadcast:  make(chan *SyncBroadcast),
			clients:    make(map[*SyncClient]bool),
			Register:   make(chan *SyncClient),
			Unregister: make(chan *SyncClient),
		}
	}
	return instance
}

// Publish the publish information in websocket
func (hub *SyncHub) Publish(w http.ResponseWriter, gp *core.Goploy) {
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
	projectUsers, err := model.ProjectUser{UserID: gp.TokenInfo.ID}.GetListByUserID()
	if err != nil || len(projectUsers) == 0 {
		c.WriteJSON(&SyncBroadcast{
			DataType: 0,
			Message:  "没有绑定服务器",
		})
		c.Close()
		return
	}
	projectMap := make(map[uint32]struct{})
	for _, projectUser := range projectUsers {
		projectMap[projectUser.ProjectID] = struct{}{}
	}
	hub.Register <- &SyncClient{
		Conn:       c,
		UserID:     gp.TokenInfo.ID,
		UserName:   gp.TokenInfo.Name,
		ProjectMap: projectMap,
	}
}

// Run goroutine run the sync hub
func (hub *SyncHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Conn.Close()
			}
		case broadcast := <-hub.Broadcast:
			for client := range hub.clients {
				if client.UserID != broadcast.UserID {
					continue
				}
				if _, ok := client.ProjectMap[broadcast.ProjectID]; ok {
					if err := client.Conn.WriteJSON(broadcast); websocket.IsCloseError(err) {
						hub.Unregister <- client
					}
				}
			}
		}
	}
}
