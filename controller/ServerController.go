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
type Server Controller

// GetList server list
func (server Server) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Server model.Servers `json:"serverList"`
	}

	serverList, err := model.Server{}.GetList()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Server: serverList}}
	response.JSON(w)
}

// GetOption server list
func (server Server) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		Server model.Servers `json:"serverList"`
	}

	serverList, err := model.Server{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RepData{Server: serverList}}
	response.JSON(w)
}

// Add one server
func (server Server) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name  string `json:"name"`
		IP    string `json:"ip"`
		Owner string `json:"owner"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err = model.Server{
		Name:       reqData.Name,
		IP:         reqData.IP,
		Owner:      reqData.Owner,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one server
func (server Server) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID    uint32 `json:"id"`
		Name  string `json:"name"`
		IP    string `json:"ip"`
		Owner string `json:"owner"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(gp.Request.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Server{
		ID:         reqData.ID,
		Name:       reqData.Name,
		IP:         reqData.IP,
		Owner:      reqData.Owner,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}
