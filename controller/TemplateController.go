package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"goploy/core"
	"goploy/model"
)

// Template struct
type Template Controller

// GetList template list
func (template Template) GetList(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Template   model.Templates  `json:"templateList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	templateList, pagination, err := model.Template{}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Template: templateList, Pagination: pagination}}
}

// GetOption template list
func (template Template) GetOption(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Template model.Templates `json:"templateList"`
	}

	templateList, err := model.Template{}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Template: templateList}}
}

// Add one template
func (template Template) Add(w http.ResponseWriter, gp *core.Goploy) *core.Response {
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
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit one template
func (template Template) Edit(w http.ResponseWriter, gp *core.Goploy) *core.Response {
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
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Remove one Template
func (template Template) Remove(w http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err = model.Template{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "删除成功"}
}
