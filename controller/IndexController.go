package controller

import (
	"net/http"

	"goploy/core"
	"goploy/model"
)

// Index struct
type Index Controller

// Get user list
func (index Index) Get(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Charts model.Charts `json:"charts"`
	}
	model := model.Charts{}
	date := gp.URLQuery.Get("date")
	if len(date) != 10 {
		response := core.Response{Code: core.Deny, Message: "日期参数错误"}
		response.JSON(w)
		return
	}
	err := model.Query(date)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Data: RespData{Charts: model}}
	response.JSON(w)
}
