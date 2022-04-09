// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"strconv"
)

func AddUploadLog(gp *core.Goploy, resp core.Response) {
	var serverID int64 = 0
	var path = ""
	respJson := resp.(response.JSON)
	if respJson.Code != response.IllegalParam {
		serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		file, fileHandler, _ := gp.Request.FormFile("file")
		path = gp.URLQuery.Get("filePath") + "/" + fileHandler.Filename
		_ = file.Close()
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypeUpload,
		Path:        path,
		Reason:      respJson.Message,
	}.AddRow()
	if err != nil {
		core.Log(core.ERROR, err.Error())
	}
}

func AddDownloadLog(gp *core.Goploy, resp core.Response) {
	msg := ""
	path := ""
	var serverID int64 = 0
	switch resp.(type) {
	case response.JSON:
		respJson := resp.(response.JSON)
		if respJson.Code != response.IllegalParam {
			msg = respJson.Message
			serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		}
	case response.SftpFile:
		path = resp.(response.SftpFile).Filename
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypeDownload,
		Path:        path,
		Reason:      msg,
	}.AddRow()
	if err != nil {
		core.Log(core.ERROR, err.Error())
	}
}

func AddPreviewLog(gp *core.Goploy, resp core.Response) {
	msg := ""
	path := ""
	var serverID int64 = 0
	switch resp.(type) {
	case response.JSON:
		respJson := resp.(response.JSON)
		if respJson.Code != response.IllegalParam {
			msg = respJson.Message
			serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		}
	case response.SftpFile:
		path = resp.(response.SftpFile).Filename
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypePreview,
		Path:        path,
		Reason:      msg,
	}.AddRow()
	if err != nil {
		core.Log(core.ERROR, err.Error())
	}
}
