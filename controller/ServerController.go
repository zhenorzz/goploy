package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/utils"
)

// Server struct
type Server Controller

// GetList -
func (Server) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// GetTotal -
func (Server) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Server{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetOption -
func (Server) GetOption(gp *core.Goploy) *core.Response {
	serverList, err := model.Server{NamespaceID: gp.Namespace.ID}.GetAll()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Servers model.Servers `json:"list"`
		}{Servers: serverList},
	}
}

// Check server
func (Server) Check(gp *core.Goploy) *core.Response {
	type ReqData struct {
		IP       string `json:"ip" validate:"required,ip|hostname"`
		Port     int    `json:"port" validate:"min=0,max=65535"`
		Owner    string `json:"owner" validate:"required,max=255"`
		Path     string `json:"path" validate:"required,max=255"`
		Password string `json:"password" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if _, err := utils.ConnectSSH(reqData.Owner, reqData.Password, reqData.Path, reqData.IP, reqData.Port); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "Connected"}
}

// Add server
func (Server) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"required,ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password"`
		Description string `json:"description" validate:"max=255"`
	}

	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Server{
		NamespaceID: gp.Namespace.ID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
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

// Edit server
func (Server) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID          int64  `json:"id" validate:"gt=0"`
		Name        string `json:"name" validate:"required"`
		IP          string `json:"ip" validate:"required,ip|hostname"`
		Port        int    `json:"port" validate:"min=0,max=65535"`
		Owner       string `json:"owner" validate:"required,max=255"`
		Path        string `json:"path" validate:"required,max=255"`
		Password    string `json:"password" validate:"max=255"`
		Description string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Server{
		ID:          reqData.ID,
		Name:        reqData.Name,
		IP:          reqData.IP,
		Port:        reqData.Port,
		Owner:       reqData.Owner,
		Path:        reqData.Path,
		Password:    reqData.Password,
		Description: reqData.Description,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoveRow server
func (Server) Remove(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Server{ID: reqData.ID}).RemoveRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}
