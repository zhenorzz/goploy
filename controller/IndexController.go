package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Index struct
type Index struct{}

// Get user list
func (index *Index) Get(w http.ResponseWriter, r *http.Request) {
	type RespData struct {
		Charts model.Charts `json:"charts"`
	}
	model := model.Charts{}
	date := r.URL.Query().Get("date")
	if len(date) != 10 {
		response := core.Response{Code: 1, Message: "日期参数错误"}
		response.Json(w)
		return
	}
	err := model.Query(date)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	response := core.Response{Data: RespData{Charts: model}}
	response.Json(w)
}

type msg struct {
	Num int
}

// Echo user list
func (index *Index) Echo(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if strings.Contains(r.Header.Get("origin"), strings.Split(r.Host, ":")[0]) {
				return true
			} else {
				return false
			}
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	m := msg{1}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteJSON(m)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
