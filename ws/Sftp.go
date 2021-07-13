package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Sftp the server file information in websocket
func (hub *Hub) Sftp(gp *core.Goploy) *core.Response {
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer c.Close()
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return nil
	}
	server, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return nil
	}
	client, err := utils.DialSSH(server.Owner, server.Password, server.Path, server.IP, server.Port)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return nil
	}
	defer client.Close()

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return nil
	}

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
				core.Log(core.ERROR, err.Error())
			}
			break
		}
		if messageType == websocket.TextMessage {
			var fileList []fileInfo
			code := core.Pass
			msg := ""
			fileInfos, err := sftpClient.ReadDir(string(message))
			if err != nil {
				code = core.Error
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

			b, _ := json.Marshal(core.Response{Code: code, Message: msg, Data: fileList})
			if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
				core.Log(core.ERROR, err.Error())
				break
			}
		}
	}

	return nil
}
