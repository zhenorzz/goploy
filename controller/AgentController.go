package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"net/http"
)

type Agent Controller

func (a Agent) Routes() []core.Route {
	return []core.Route{
		core.NewWhiteRoute("/agent/report", http.MethodPost, a.Report),
	}
}

func (Agent) Report(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerId   int64  `json:"serverId" validate:"gt=0"`
		Type       int    `json:"type" validate:"gt=0"`
		Item       string `json:"item" validate:"required"`
		Value      string `json:"value" validate:"required"`
		ReportTime string `json:"reportTime" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ServerAgentLog{
		ServerID:   reqData.ServerId,
		Type:       reqData.Type,
		Item:       reqData.Item,
		Value:      reqData.Value,
		ReportTime: reqData.ReportTime,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
