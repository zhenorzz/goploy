package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Server struct
type Server struct{}

// Get server list
func (server *Server) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Server model.Servers `json:"serverList"`
	}

	model := model.Servers{}
	err := model.Query()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{Server: model}}
	response.Json(w)
}

// Add one server
func (server *Server) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		Path string `json:"path"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.Server{
		Name:       reqData.Name,
		IP:         reqData.IP,
		Path:       reqData.Path,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = model.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.Json(w)
}
