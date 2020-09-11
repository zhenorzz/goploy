package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"net"
	"strconv"
	"time"
)

// Monitor struct
type Monitor Controller

// GetList monitor list
func (monitor Monitor) GetList(gp *core.Goploy) *core.Response {
	type RespData struct {
		Monitors model.Monitors `json:"list"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	monitorList, err := model.Monitor{NamespaceID: gp.Namespace.ID}.GetList(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Monitors: monitorList}}
}

// GetList monitor list
func (monitor Monitor) GetTotal(gp *core.Goploy) *core.Response {
	type RespData struct {
		Total int64 `json:"total"`
	}
	total, err := model.Monitor{NamespaceID: gp.Namespace.ID}.GetTotal()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Total: total}}
}

// Check one monitor
func (monitor Monitor) Check(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Domain string `json:"domain" validate:"required"`
		Port   int    `json:"port" validate:"min=0,max=65535"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	_, err := net.DialTimeout("tcp", reqData.Domain+":"+strconv.Itoa(reqData.Port), 5*time.Second)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Message: "Connected"}
}

// Add one monitor
func (monitor Monitor) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name         string `json:"name" validate:"required"`
		Domain       string `json:"domain" validate:"required"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		NotifyTimes  uint16 `json:"notifyTimes" validate:"gt=0"`
		Description  string `json:"description" validate:"max=255"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := model.Monitor{
		NamespaceID:  gp.Namespace.ID,
		Name:         reqData.Name,
		Domain:       reqData.Domain,
		Port:         reqData.Port,
		Second:       reqData.Second,
		Times:        reqData.Times,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
		NotifyTimes:  reqData.NotifyTimes,
		Description:  reqData.Description,
	}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{ID: id}}
}

// Edit one monitor
func (monitor Monitor) Edit(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		Name         string `json:"name" validate:"required"`
		Domain       string `json:"domain" validate:"required"`
		Port         int    `json:"port" validate:"min=0,max=65535"`
		Second       int    `json:"second" validate:"gt=0"`
		Times        uint16 `json:"times" validate:"gt=0"`
		NotifyType   uint8  `json:"notifyType" validate:"gt=0"`
		NotifyTarget string `json:"notifyTarget" validate:"required"`
		NotifyTimes  uint16 `json:"notifyTimes" validate:"gt=0"`
		Description  string `json:"description" validate:"max=255"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Monitor{
		ID:           reqData.ID,
		Name:         reqData.Name,
		Domain:       reqData.Domain,
		Port:         reqData.Port,
		Second:       reqData.Second,
		Times:        reqData.Times,
		NotifyType:   reqData.NotifyType,
		NotifyTarget: reqData.NotifyTarget,
		NotifyTimes:  reqData.NotifyTimes,
		Description:  reqData.Description,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// Stop one Monitor
func (monitor Monitor) Toggle(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID}).ToggleState(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// RemoveRow one Monitor
func (monitor Monitor) Remove(gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.Monitor{ID: reqData.ID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}
