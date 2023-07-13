// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"golang.org/x/crypto/ssh"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// write data to WebSocket
// the data comes from ssh server.
type xtermBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *xtermBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func (hub *Hub) xterm(gp *server.Goploy) server.Response {
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
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	defer c.Close()
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	rows, err := strconv.Atoi(gp.URLQuery.Get("rows"))
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	cols, err := strconv.Atoi(gp.URLQuery.Get("cols"))
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}

	srv, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	defer client.Close()
	// create session
	session, err := client.NewSession()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	defer session.Close()
	sessionStdin, err := session.StdinPipe()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	comboWriter := new(xtermBufferWriter)
	//ssh.stdout and stderr will write output into comboWriter
	session.Stdout = comboWriter
	session.Stderr = comboWriter
	// Request pseudo terminal
	if err := session.RequestPty("xterm", rows, cols, ssh.TerminalModes{}); err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}

	// terminal log
	tlID, err := model.TerminalLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		StartTime:   time.Now().Format("20060102150405"),
	}.AddRow()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}

	var recorder *pkg.Recorder
	recorder, err = pkg.NewRecorder(config.GetTerminalLogPath(tlID), "xterm", rows, cols)
	if err != nil {
		log.Error(err.Error())
	} else {
		defer recorder.Close()
	}

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	flushMessageTick := time.NewTicker(time.Millisecond * time.Duration(50))
	defer flushMessageTick.Stop()
	stop := make(chan bool, 1)
	defer func() {
		stop <- true
	}()
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					c.Close()
					return
				}
			case <-flushMessageTick.C:
				if comboWriter.buffer.Len() != 0 {
					err := c.WriteMessage(websocket.BinaryMessage, comboWriter.buffer.Bytes())
					if err != nil {
						c.Close()
						return
					}
					if recorder != nil {
						if err := recorder.WriteData(comboWriter.buffer.String()); err != nil {
							log.Error(err.Error())
						}
					}
					comboWriter.buffer.Reset()
				}
			case <-stop:
				return
			}
		}
	}()

	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error(err.Error())
			}
			break
		}
		if messageType != websocket.PongMessage {
			if _, err := sessionStdin.Write(p); err != nil {
				log.Error(err.Error())
				break
			}
		}
	}

	if err := (model.TerminalLog{
		ID:      tlID,
		EndTime: time.Now().Format("20060102150405"),
	}.EditRow()); err != nil {
		log.Error(err.Error())
	}

	return response.Empty{}
}
