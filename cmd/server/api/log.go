// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/zhenorzz/goploy/config"
	model2 "github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

// Log struct
type Log API

func (l Log) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/log/getLoginLogList", http.MethodGet, l.GetLoginLogList).Permissions(config.ShowLoginLogPage),
		server.NewRoute("/log/getLoginLogTotal", http.MethodGet, l.GetLoginLogTotal).Permissions(config.ShowLoginLogPage),
		server.NewRoute("/log/getOperationLogList", http.MethodGet, l.GetOperationLogList).Permissions(config.ShowOperationLogPage),
		server.NewRoute("/log/getOperationLogTotal", http.MethodGet, l.GetOperationLogTotal).Permissions(config.ShowOperationLogPage),
		server.NewRoute("/log/getSftpLogList", http.MethodGet, l.GetSftpLogList).Permissions(config.ShowSFTPLogPage),
		server.NewRoute("/log/getSftpLogTotal", http.MethodGet, l.GetSftpLogTotal).Permissions(config.ShowSFTPLogPage),
		server.NewRoute("/log/getTerminalLogList", http.MethodGet, l.GetTerminalLogList).Permissions(config.ShowTerminalLogPage),
		server.NewRoute("/log/getTerminalLogTotal", http.MethodGet, l.GetTerminalLogTotal).Permissions(config.ShowTerminalLogPage),
		server.NewRoute("/log/getTerminalRecord", http.MethodGet, l.GetTerminalRecord).Permissions(config.ShowTerminalRecord),
		server.NewRoute("/log/getPublishLogList", http.MethodGet, l.GetPublishLogList).Permissions(config.ShowPublishLogPage),
		server.NewRoute("/log/getPublishLogTotal", http.MethodGet, l.GetPublishLogTotal).Permissions(config.ShowPublishLogPage),
	}
}

func (Log) GetLoginLogList(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account string `schema:"account"`
		Page    uint64 `schema:"page" validate:"gt=0"`
		Rows    uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	list, err := model2.LoginLog{Account: reqData.Account}.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model2.LoginLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetLoginLogTotal(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account string `schema:"account"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	total, err := model2.LoginLog{Account: reqData.Account}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Log) GetOperationLogList(gp *server.Goploy) server.Response {
	type ReqData struct {
		Username string `schema:"username"`
		Router   string `schema:"router"`
		API      string `schema:"api"`
		Page     uint64 `schema:"page" validate:"gt=0"`
		Rows     uint64 `schema:"rows" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	opLog := model2.OperationLog{Username: reqData.Username, Router: reqData.Router, API: reqData.API}

	if gp.UserInfo.SuperManager != model2.SuperManager {
		opLog.NamespaceID = gp.Namespace.ID
	}

	list, err := opLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model2.OperationLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetOperationLogTotal(gp *server.Goploy) server.Response {
	type ReqData struct {
		Username string `schema:"username"`
		Router   string `schema:"router"`
		API      string `schema:"api"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	opLog := model2.OperationLog{Username: reqData.Username, Router: reqData.Router, API: reqData.API}

	if gp.UserInfo.SuperManager != model2.SuperManager {
		opLog.NamespaceID = gp.Namespace.ID
	}

	total, err := opLog.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (Log) GetSftpLogList(gp *server.Goploy) server.Response {
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

	sftpLog := model2.SftpLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
		sftpLog.NamespaceID = gp.Namespace.ID
	}

	list, err := sftpLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model2.SftpLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetSftpLogTotal(gp *server.Goploy) server.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	sftpLog := model2.SftpLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
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

func (Log) GetTerminalLogList(gp *server.Goploy) server.Response {
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

	terminalLog := model2.TerminalLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
		terminalLog.NamespaceID = gp.Namespace.ID
	}

	list, err := terminalLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model2.TerminalLogs `json:"list"`
		}{List: list},
	}
}

func (Log) GetTerminalLogTotal(gp *server.Goploy) server.Response {
	type ReqData struct {
		Username   string `schema:"username"`
		ServerName string `schema:"serverName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	terminalLog := model2.TerminalLog{Username: reqData.Username, ServerName: reqData.ServerName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
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

func (Log) GetTerminalRecord(gp *server.Goploy) server.Response {
	type ReqData struct {
		RecordID int64 `schema:"recordId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}
	terminalLog, err := model2.TerminalLog{ID: reqData.RecordID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if gp.UserInfo.SuperManager != model2.SuperManager && terminalLog.NamespaceID != gp.Namespace.ID {
		return response.JSON{Code: response.Error, Message: "You have no access to enter this record"}
	}
	return response.File{Filename: config.GetTerminalLogPath(reqData.RecordID)}
}

func (Log) GetPublishLogList(gp *server.Goploy) server.Response {
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

	publishLog := model2.PublishTrace{PublisherName: reqData.Username, ProjectName: reqData.ProjectName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
		publishLog.NamespaceID = gp.Namespace.ID
	}

	list, err := publishLog.GetList(reqData.Page, reqData.Rows)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model2.PublishTraces `json:"list"`
		}{List: list},
	}
}

func (Log) GetPublishLogTotal(gp *server.Goploy) server.Response {
	type ReqData struct {
		Username    string `schema:"username"`
		ProjectName string `schema:"projectName"`
	}
	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	publishLog := model2.PublishTrace{PublisherName: reqData.Username, ProjectName: reqData.ProjectName}

	if gp.UserInfo.SuperManager != model2.SuperManager {
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
