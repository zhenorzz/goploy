package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// Cron struct
type Cron Controller

// Report -
func (Cron) Report(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ServerId   int64  `json:"serverId" validate:"gt=0"`
		CronId     int64  `json:"cronId" validate:"gt=0"`
		ExecCode   int    `json:"execCode"`
		Message    string `json:"message" validate:"required"`
		ReportTime string `json:"reportTime" validate:"required"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err := model.CronLog{
		ServerID:   reqData.ServerId,
		CronID:     reqData.CronId,
		ExecCode:   reqData.ExecCode,
		Message:    reqData.Message,
		ReportTime: reqData.ReportTime,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// GetList cron list
func (Cron) GetList(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ServerID int64 `json:"serverId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crons, err := model.Cron{ServerID: reqData.ServerID}.GetList()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			List model.Crons `json:"list"`
		}{List: crons},
	}
}

// GetLogs cron log list
func (Cron) GetLogs(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	type ReqData struct {
		ServerID int64 `json:"serverId" validate:"gt=0"`
		CronID   int64 `json:"cronId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	crons, err := model.CronLog{ServerID: reqData.ServerID, CronID: reqData.CronID}.GetList(pagination)

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			List model.CronLogs `json:"list"`
		}{List: crons},
	}
}

// Add cron
func (Cron) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ServerID    int64  `json:"serverId" validate:"gt=0"`
		Expression  string `json:"expression" validate:"required"`
		Command     string `json:"command" validate:"required"`
		SingleMode  uint8  `json:"singleMode" validate:"gte=0"`
		LogLevel    uint8  `json:"logLevel" validate:"gte=0"`
		Description string `json:"description" validate:"max=255"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}

	}
	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// Edit cron
func (Cron) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Expression  string `json:"expression" validate:"required"`
		Command     string `json:"command" validate:"required"`
		SingleMode  uint8  `json:"singleMode" validate:"gte=0"`
		LogLevel    uint8  `json:"logLevel" validate:"gte=0"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
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
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Remove cron
func (Cron) Remove(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Cron{ID: reqData.ID}).RemoveRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{}
}
