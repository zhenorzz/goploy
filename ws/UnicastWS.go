package ws

// 单播
import (
	"log"
	"net/http"
	"strings"

	"goploy/core"

	"github.com/gorilla/websocket"
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

var unicastHub *UnicastHub

// GetUnicastHub it will only init once in main.go
func GetUnicastHub() *UnicastHub {
	if unicastHub == nil {
		unicastHub = &UnicastHub{
			UnicastData: make(chan *UnicastData),
			clients:     make(map[*UnicastClient]bool),
			Register:    make(chan *UnicastClient),
			Unregister:  make(chan *UnicastClient),
		}
	}
	return unicastHub
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
