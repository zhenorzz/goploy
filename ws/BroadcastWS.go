package ws

// 单播
import (
	"goploy/model"
	"net/http"
	"strings"
	"time"

	"goploy/core"

	"github.com/gorilla/websocket"
)

// BroadcastClient stores a client information
type BroadcastClient struct {
	Conn     *websocket.Conn
	UserInfo model.User
}

// BroadcastData is message struct
type BroadcastData struct {
	Type    int
	Message interface{}
}

// ProjectMessage is publish project message struct
type ProjectMessage struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	State       uint8  `json:"state"`
	Message     string `json:"message"`
}

// BroadcastHub is a client struct
type BroadcastHub struct {
	// Registered clients.
	clients map[*BroadcastClient]bool

	// Inbound messages from the clients.
	BroadcastData chan *BroadcastData

	// Register requests from the clients.
	Register chan *BroadcastClient

	// Unregister requests from clients.
	Unregister chan *BroadcastClient
	// ping pong ticker
	ticker chan *BroadcastClient
}

const (
	TypeProject = 1
)

var broadcastHub *BroadcastHub

// GetBroadcastHub it will only init once in main.go
func GetBroadcastHub() *BroadcastHub {
	if broadcastHub == nil {
		broadcastHub = &BroadcastHub{
			BroadcastData: make(chan *BroadcastData),
			clients:       make(map[*BroadcastClient]bool),
			Register:      make(chan *BroadcastClient),
			Unregister:    make(chan *BroadcastClient),
			ticker:        make(chan *BroadcastClient),
		}
	}
	return broadcastHub
}

// Broadcast the publish information in websocket
func (hub *BroadcastHub) Broadcast(w http.ResponseWriter, gp *core.Goploy) {
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
		core.Log(core.ERROR, err.Error())
		return
	}
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	broadcastClient := &BroadcastClient{
		Conn:     c,
		UserInfo: gp.UserInfo,
	}
	hub.Register <- broadcastClient

	ticker := time.NewTicker(pingPeriod)
	stop := make(chan bool, 1)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				hub.ticker <- broadcastClient
			case <-stop:
				return
			}
		}
	}(ticker)
	// you must read message to trigger pong handler
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				core.Log(core.ERROR, err.Error())
			}
			break
		}
	}

	defer func() {
		hub.Unregister <- broadcastClient
		c.Close()
		ticker.Stop()
		stop <- true
	}()
}

// Run goroutine run the sync hub
func (hub *BroadcastHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Conn.Close()
			}
		case broadcast := <-hub.BroadcastData:
			for client := range hub.clients {
				if broadcast.Type == TypeProject {
					projectMessage := broadcast.Message.(ProjectMessage)
					_, err := model.Project{ID: projectMessage.ProjectID}.GetUserProjectData(client.UserInfo.ID, client.UserInfo.Role, client.UserInfo.ManageGroupStr)
					if err != nil {
						continue
					}
				} else {
					continue
				}
				if err := client.Conn.WriteJSON(broadcast.Message); websocket.IsCloseError(err) {
					hub.Unregister <- client
				}
			}
		case client := <-hub.ticker:
			if _, ok := hub.clients[client]; ok {
				if err := client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					hub.Unregister <- client
				}
			}
		}
	}
}
