package controller

import (
	"net/http"

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
