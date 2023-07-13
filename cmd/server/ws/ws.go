// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
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
	maxMessageSize = 10240
)

const (
	TypeProject = 1
	TypeMonitor = 3
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
	CanSendTo(client *Client) error
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

func init() {
	go hub.run()
}

func (hub *Hub) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/ws/connect", http.MethodGet, hub.connect),
		server.NewRoute("/ws/xterm", http.MethodGet, hub.xterm),
		server.NewRoute("/ws/sftp", http.MethodGet, hub.sftp),
	}
}

var hub = &Hub{
	Data:       make(chan *Data),
	clients:    make(map[*Client]bool),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	ticker:     make(chan *Client),
}

func GetHub() *Hub {
	return hub
}

func Send(d Data) {
	GetHub().Data <- &d
}

func (hub *Hub) connect(gp *server.Goploy) server.Response {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if config.Toml.CORS.Enabled {
				if config.Toml.CORS.Origins == "*" {
					return true
				} else if strings.Contains(config.Toml.CORS.Origins, r.Header.Get("origin")) {
					return true
				}
			}
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			}
			return false
		},
	}
	c, err := upgrader.Upgrade(gp.ResponseWriter, gp.Request, nil)
	if err != nil {
		log.Error(err.Error())
		return response.JSON{Code: response.Error, Message: err.Error()}
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
	go func() {
		for {
			select {
			case <-ticker.C:
				hub.ticker <- client
			case <-stop:
				return
			}
		}
	}()
	// you must read message to trigger pong handler
	for {
		_, _, err = c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error(err.Error())
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

	return response.Empty{}
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
				if data.Message.CanSendTo(client) != nil {
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
