package ws

// 单播
import (
	"log"
	"net/http"
	"strings"

	"goploy/core"

	"github.com/gorilla/websocket"
)

// BroadcastClient stores a client information
type BroadcastClient struct {
	Conn       *websocket.Conn
	UserID     int64
	UserName   string
}

// BroadcastData is message struct
type BroadcastData struct {
	Type    int
	Message interface{}
}

// ProjectMessage is publish project message struct
type ProjectMessage1 struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	UserID      int64  `json:"userId"`
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
		log.Print("upgrade:", err)
		return
	}
	hub.Register <- &BroadcastClient{
		Conn:     c,
		UserID:   gp.TokenInfo.ID,
		UserName: gp.TokenInfo.Name,
	}
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
					if ok := core.UserHasProject(client.UserID, projectMessage.ProjectID); !ok {
						continue
					}
				} else {
					continue
				}

				if err := client.Conn.WriteJSON(broadcast.Message); websocket.IsCloseError(err) {
					hub.Unregister <- client
				}
			}
		}
	}
}
