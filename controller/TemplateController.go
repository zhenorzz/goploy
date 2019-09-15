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
func (template Template) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Template   model.Templates  `json:"templateList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	templateList, pagination, err := model.Template{}.GetList(pagination)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{Template: templateList, Pagination: pagination}}
	response.JSON(w)
}

// GetOption template list
func (template Template) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		Template model.Templates `json:"templateList"`
	}

	templateList, err := model.Template{}.GetAll()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{Template: templateList}}
	response.JSON(w)
}

// Add one template
func (template Template) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Name         string `json:"name"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script"`
	}
	var reqData ReqData

	if err := json.Unmarshal(gp.Body, &reqData); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	_, err := model.Template{
		Name:         reqData.Name,
		Remark:       reqData.Remark,
		PackageIDStr: reqData.PackageIDStr,
		Script:       reqData.Script,
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one template
func (template Template) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Template{
		ID:           reqData.ID,
		Name:         reqData.Name,
		Remark:       reqData.Remark,
		PackageIDStr: reqData.PackageIDStr,
		Script:       reqData.Script,
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}

// Remove one Template
func (template Template) Remove(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.Template{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.Remove()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}
