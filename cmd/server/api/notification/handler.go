// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package notification

import (
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
)

type Notification api.API

func (n Notification) Handler() []server.Route {
	return []server.Route{
		server.NewRoute("/notification/getList", http.MethodGet, n.GetList).Permissions(config.ShowNotificationPage),
		server.NewRoute("/notification/edit", http.MethodPut, n.Edit).Permissions(config.EditNotification).LogFunc(middleware.AddOPLog),
	}
}

// GetList lists notification template
// @Summary List notification template
// @Tags Notification
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=notification.GetList.RespData}
// @Router /notification/getList [get]
func (Notification) GetList(*server.Goploy) server.Response {
	list, err := model.NotificationTemplate{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		List model.NotificationTemplates `json:"list"`
	}

	return response.JSON{
		Data: RespData{List: list},
	}
}

// Edit edits the notification template
// @Summary Edit the notification template
// @Tags Notification
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body notification.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /notification/edit [put]
func (Notification) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID       int64  `json:"id" validate:"required,gt=0"`
		Title    string `json:"title" validate:"required"`
		Template string `json:"template" validate:"required"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.NotificationTemplate{ID: reqData.ID, Title: reqData.Title, Template: reqData.Template}.EditRow()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}
