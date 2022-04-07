package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/response"
	"net/http"
)

// Log struct
type Log Controller

func (l Log) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/log/getLoginLogList", http.MethodGet, l.GetLoginLogList).Permissions(permission.ShowLoginLogPage),
		core.NewRoute("/log/getLoginLogTotal", http.MethodGet, l.GetLoginLogTotal).Permissions(permission.ShowLoginLogPage),
		core.NewRoute("/log/getSftpLogList", http.MethodGet, l.GetSftpLogList).Permissions(permission.ShowSFTPLogPage),
		core.NewRoute("/log/getSftpLogTotal", http.MethodGet, l.GetSftpLogTotal).Permissions(permission.ShowSFTPLogPage),
		core.NewRoute("/log/getTerminalLogList", http.MethodGet, l.GetTerminalLogList).Permissions(permission.ShowTerminalLogPage),
		core.NewRoute("/log/getTerminalLogTotal", http.MethodGet, l.GetTerminalLogTotal).Permissions(permission.ShowTerminalLogPage),
		core.NewRoute("/log/getTerminalRecord", http.MethodGet, l.GetTerminalRecord).Permissions(permission.ShowTerminalRecord),
		core.NewRoute("/log/getPublishLogList", http.MethodGet, l.GetPublishLogList).Permissions(permission.ShowPublishLogPage),
		core.NewRoute("/log/getPublishLogTotal", http.MethodGet, l.GetPublishLogTotal).Permissions(permission.ShowPublishLogPage),
	}
}

func (Log) GetLoginLogList(gp *core.Goploy) core.Response {
	type ReqData struct {
		Account string `schema:"account"`
		Page    uint64 `schema:"page" validate:"gt=0"`
		Rows    uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	list, err := model.LoginLog{Account: reqData.Account}.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.LoginLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetLoginLogTotal(gp *core.Goploy) core.Response {
	type ReqData struct {
		Account string `schema:"account"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	total, err := model.LoginLog{Account: reqData.Account}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Log) GetSftpLogList(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
		Page       uint64 `schema:"page" validate:"gt=0"`
		Rows       uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	sftpLog := model.SftpLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		sftpLog.NamespaceID = gp.Namespace.ID
	}

	list, err := sftpLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.SftpLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetSftpLogTotal(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	sftpLog := model.SftpLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		sftpLog.NamespaceID = gp.Namespace.ID
	}

	total, err := sftpLog.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Log) GetTerminalLogList(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
		Page       uint64 `schema:"page" validate:"gt=0"`
		Rows       uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	terminalLog := model.TerminalLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		terminalLog.NamespaceID = gp.Namespace.ID
	}

	list, err := terminalLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.TerminalLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetTerminalLogTotal(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	terminalLog := model.TerminalLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		terminalLog.NamespaceID = gp.Namespace.ID
	}

	total, err := terminalLog.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Log) GetTerminalRecord(gp *core.Goploy) core.Response {
	type ReqData struct {
		RecordID int64 `schema:"recordId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	terminalLog, err := model.TerminalLog{ID: reqData.RecordID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if gp.UserInfo.SuperManager != model.SuperManager && terminalLog.NamespaceID != gp.Namespace.ID {
		return response.JSON{Code: response.Error, Message: "You have no access to enter this record"}
	}
	return response.File{Filename: core.GetTerminalLogPath(reqData.RecordID)}
}

func (Log) GetPublishLogList(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username    string `schema:"username"`
		ProjectName string `schema:"projectName"`
		Page        uint64 `schema:"page" validate:"gt=0"`
		Rows        uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	publishLog := model.PublishTrace{PublisherName: reqData.Username, ProjectName: reqData.ProjectName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		publishLog.NamespaceID = gp.Namespace.ID
	}

	list, err := publishLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.PublishTraces `json:"list"`
		}{List: list},
	}
}

func (Log) GetPublishLogTotal(gp *core.Goploy) core.Response {
	type ReqData struct {
		Username    string `schema:"username"`
		ProjectName string `schema:"projectName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	publishLog := model.PublishTrace{PublisherName: reqData.Username, ProjectName: reqData.ProjectName}

	if gp.UserInfo.SuperManager != model.SuperManager {
		publishLog.NamespaceID = gp.Namespace.ID
	}

	total, err := publishLog.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}
