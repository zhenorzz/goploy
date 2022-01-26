package ws

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/utils"
	"golang.org/x/crypto/ssh"
	"net/http"
	"path"
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
func (hub *Hub) Xterm(gp *core.Goploy) core.Response {
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

	server, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		return response.Empty{}
	}
	client, err := server.Convert2SSHConfig().Dial()
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

	recorder, _ := utils.NewRecorder(path.Join(core.GetLogPath(), "demo.cast"), "xterm", rows, cols)
	defer recorder.Close()
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
					recorder.WriteData(comboWriter.buffer.String())
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

	return response.Empty{}
}
