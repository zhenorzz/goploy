package ws

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
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

// Xterm the publish information in websocket
func (hub *Hub) Xterm(gp *core.Goploy) *core.Response {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			}
			return false
		},
	}
	rows, err := strconv.Atoi(gp.URLQuery.Get("rows"))
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	cols, err := strconv.Atoi(gp.URLQuery.Get("cols"))
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverID, err := strconv.ParseInt(gp.URLQuery.Get("serverId"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	c, err := upgrader.Upgrade(gp.ResponseWriter, gp.Request, nil)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer c.Close()
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	server, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	client, err := utils.DialSSH(server.Owner, server.Password, server.Path, server.IP, server.Port)
	if err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer client.Close()
	// create session
	session, err := client.NewSession()
	if err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	defer session.Close()
	sessionStdin, err := session.StdinPipe()
	if err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	comboWriter := new(xtermBufferWriter)
	//ssh.stdout and stderr will write output into comboWriter
	session.Stdout = comboWriter
	session.Stderr = comboWriter
	// Request pseudo terminal
	if err := session.RequestPty("xterm", rows, cols, ssh.TerminalModes{}); err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		core.Log(core.ERROR, err.Error())
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	flushMessageTick := time.NewTicker(time.Millisecond * time.Duration(120))
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
				core.Log(core.ERROR, err.Error())
			}
			break
		}
		if messageType != websocket.PongMessage {
			if _, err := sessionStdin.Write(p); err != nil {
				core.Log(core.ERROR, err.Error())
				break
			}
		}
	}

	return nil
}
