package ws

import (
	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"net/http"
	"strings"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

const (
	TypeProject        = 1
	TypeServerTemplate = 2
	TypeMonitor        = 3
)

// Client stores a client information
type Client struct {
	Conn     *websocket.Conn
	UserInfo model.User
}

// Data is message struct
type Data struct {
	Type    int
	UserIDs []int64
	Message Message
}

type Message interface {
	canSendTo(client *Client) error
}

// Hub is a client struct
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Data chan *Data

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
	// ping pong ticker
	ticker chan *Client
}

// Init websocket
func Init() {
	go GetHub().run()
}

var hub *Hub

// GetHub it will only init once in main.go
func GetHub() *Hub {
	if hub == nil {
		hub = &Hub{
			Data:       make(chan *Data),
			clients:    make(map[*Client]bool),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
			ticker:     make(chan *Client),
		}
	}
	return hub
}

// Connect the publish information in websocket
func (hub *Hub) Connect(gp *core.Goploy) *core.Response {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			}
			return false
		},
	}
	c, err := upgrader.Upgrade(gp.ResponseWriter, gp.Request, nil)
	if err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	client := &Client{
		Conn:     c,
		UserInfo: gp.UserInfo,
	}
	hub.Register <- client

	ticker := time.NewTicker(pingPeriod)
	stop := make(chan bool, 1)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				hub.ticker <- client
			case <-stop:
				return
			}
		}
	}(ticker)
	// you must read message to trigger pong handler
	for {
		_, _, err = c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				core.Log(core.ERROR, err.Error())
			}
			break
		}
	}

	defer func() {
		hub.Unregister <- client
		c.Close()
		ticker.Stop()
		stop <- true
	}()

	return nil
}

// Run goroutine run the sync hub
func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				client.Conn.Close()
			}
		case data := <-hub.Data:
			for client := range hub.clients {
				if data.Message.canSendTo(client) != nil {
					continue
				}
				// check userIDs list
				for _, userID := range data.UserIDs {
					if client.UserInfo.ID != userID {
						continue
					}
				}
				if err := client.Conn.WriteJSON(
					struct {
						Type    int         `json:"type"`
						Message interface{} `json:"message"`
					}{
						Type:    data.Type,
						Message: data.Message,
					}); websocket.IsCloseError(err) {
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
