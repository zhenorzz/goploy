// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// sftp the server file information in websocket
func (hub *Hub) sftp(gp *server.Goploy) server.Response {
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

	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	srv, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	client, err := srv.ToSSHConfig().Dial()
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	defer sftpClient.Close()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
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
			case <-stop:
				return
			}
		}
	}()

	type fileInfo struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Mode    string `json:"mode"`
		ModTime string `json:"modTime"`
		IsDir   bool   `json:"isDir"`
	}

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error(err.Error())
			}
			break
		}
		if messageType == websocket.TextMessage {
			var fileList []fileInfo
			code := response.Pass
			msg := ""
			fileInfos, err := sftpClient.ReadDir(string(message))
			if err != nil {
				code = response.Error
				msg = err.Error()
			} else {
				for _, f := range fileInfos {
					if f.Mode()&os.ModeSymlink != 0 {
						continue
					}
					fileList = append(fileList, fileInfo{
						Name:    f.Name(),
						Size:    f.Size(),
						Mode:    f.Mode().String(),
						ModTime: f.ModTime().Format("2006-01-02 15:04:05"),
						IsDir:   f.IsDir(),
					})
				}
			}

			b, _ := json.Marshal(response.JSON{Code: code, Message: msg, Data: fileList})
			if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Error(err.Error())
				break
			}

			// sftp log
			if err := (model.SftpLog{
				NamespaceID: gp.Namespace.ID,
				UserID:      gp.UserInfo.ID,
				ServerID:    serverID,
				RemoteAddr:  gp.Request.RemoteAddr,
				UserAgent:   gp.Request.UserAgent(),
				Type:        model.SftpLogTypeRead,
				Path:        string(message),
				Reason:      msg,
			}.AddRow()); err != nil {
				log.Error(err.Error())
			}
		}
	}

	return response.Empty{}
}
