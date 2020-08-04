package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Template struct
type Template Controller

// GetList template list
func (Template) GetList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Templates model.Templates `json:"list"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	templateList, err := model.Template{}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Templates: templateList}}
}

// GetTotal template total
func (Template) GetTotal(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Total int64 `json:"total"`
	}
	total, err := model.Template{}.GetTotal()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Total: total}}
}

// GetOption template list
func (Template) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Templates model.Templates `json:"list"`
	}

	templateList, err := model.Template{}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Templates: templateList}}
}

// Add one template
func (Template) Add(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script" validate:"required"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}

	var reqData ReqData

	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Template{
		Name:         reqData.Name,
		Remark:       reqData.Remark,
		PackageIDStr: reqData.PackageIDStr,
		Script:       reqData.Script,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit one template
func (Template) Edit(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		Name         string `json:"name" validate:"required"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Template{
		ID:           reqData.ID,
		Name:         reqData.Name,
		Remark:       reqData.Remark,
		PackageIDStr: reqData.PackageIDStr,
		Script:       reqData.Script,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// DeleteRow one Template
func (Template) Remove(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err = model.Template{
		ID: reqData.ID,
	}.DeleteRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "删除成功"}
}
