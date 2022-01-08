package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"net/http"
)

// Cron struct
type Cron Controller

func (c Cron) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/cron/report", http.MethodPost, c.Report).White(),
		core.NewRoute("/cron/getList", http.MethodPost, c.GetList).White(),
		core.NewRoute("/cron/getLogs", http.MethodPost, c.GetLogs).White(),
		core.NewRoute("/cron/add", http.MethodPost, c.Add).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/cron/edit", http.MethodPut, c.Edit).Roles(core.RoleAdmin, core.RoleManager),
		core.NewRoute("/cron/remove", http.MethodDelete, c.Remove).Roles(core.RoleAdmin, core.RoleManager),
	}
}

func (Cron) Report(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerId   int64  `json:"serverId" validate:"gt=0"`
		CronId     int64  `json:"cronId" validate:"gt=0"`
		ExecCode   int    `json:"execCode"`
		Message    string `json:"message" validate:"required"`
		ReportTime string `json:"reportTime" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.CronLog{
		ServerID:   reqData.ServerId,
		CronID:     reqData.CronId,
		ExecCode:   reqData.ExecCode,
		Message:    reqData.Message,
		ReportTime: reqData.ReportTime,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Cron) GetList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID int64 `json:"serverId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	crons, err := model.Cron{ServerID: reqData.ServerID}.GetList()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.Crons `json:"list"`
		}{List: crons},
	}
}

func (Cron) GetLogs(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type ReqData struct {
		ServerID int64 `json:"serverId" validate:"gt=0"`
		CronID   int64 `json:"cronId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	crons, err := model.CronLog{ServerID: reqData.ServerID, CronID: reqData.CronID}.GetList(pagination)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.CronLogs `json:"list"`
		}{List: crons},
	}
}

func (Cron) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID    int64  `json:"serverId" validate:"gt=0"`
		Expression  string `json:"expression" validate:"required"`
		Command     string `json:"command" validate:"required"`
		SingleMode  uint8  `json:"singleMode" validate:"gte=0"`
		LogLevel    uint8  `json:"logLevel" validate:"gte=0"`
		Description string `json:"description" validate:"max=255"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Cron{
		ServerID:    reqData.ServerID,
		Expression:  reqData.Expression,
		Command:     reqData.Command,
		SingleMode:  reqData.SingleMode,
		LogLevel:    reqData.LogLevel,
		Description: reqData.Description,
		Creator:     gp.UserInfo.Name,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}

	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Cron) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Expression  string `json:"expression" validate:"required"`
		Command     string `json:"command" validate:"required"`
		SingleMode  uint8  `json:"singleMode" validate:"gte=0"`
		LogLevel    uint8  `json:"logLevel" validate:"gte=0"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Cron{
		ID:          reqData.ID,
		Expression:  reqData.Expression,
		Command:     reqData.Command,
		SingleMode:  reqData.SingleMode,
		LogLevel:    reqData.LogLevel,
		Description: reqData.Description,
		Editor:      gp.UserInfo.Name,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Cron) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Cron{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
