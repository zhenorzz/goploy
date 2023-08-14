// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package template

import (
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Template api.API

func (t Template) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/template/getOption", http.MethodGet, t.GetOption),
		server.NewRoute("/template/add", http.MethodPost, t.Add).LogFunc(middleware.AddOPLog),
		server.NewRoute("/template/remove", http.MethodDelete, t.Remove).LogFunc(middleware.AddOPLog),
	}
}

func (Template) GetOption(gp *server.Goploy) server.Response {
	type ReqData struct {
		Type uint8 `json:"type" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	list, err := model.Template{Type: reqData.Type}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Data: struct {
		List model.Templates `json:"list"`
	}{list}}
}

func (Template) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Type        uint8  `json:"type" validate:"required"`
		Name        string `json:"name" validate:"required,max=255"`
		Content     string `json:"content" validate:"required"`
		Description string `json:"description" validate:"max=2047"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	id, err := model.Template{
		Type:        reqData.Type,
		Name:        reqData.Name,
		Content:     reqData.Content,
		Description: reqData.Description,
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

func (Template) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.Template{ID: reqData.ID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
