package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"net/http"
	"strconv"
)

// Namespace struct
type Namespace Controller

// GetList Namespace list
func (namespace Namespace) GetList(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Namespaces model.Namespaces `json:"list"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	namespaceList, err := model.Namespace{UserID: gp.UserInfo.ID}.GetListByUserID(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{Data: RespData{Namespaces: namespaceList}}
}

// GetList server list
func (namespace Namespace) GetTotal(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		Total int64 `json:"total"`
	}
	total, err := model.Namespace{UserID: gp.UserInfo.ID}.GetTotalByUserID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{Total: total}}
}

// GetBindUserList user list
func (namespace Namespace) GetUserOption(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		NamespaceUsers model.NamespaceUsers `json:"list"`
	}

	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllUserByNamespaceID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{NamespaceUsers: namespaceUsers}}
}

// GetBindUserList user list
func (namespace Namespace) GetBindUserList(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type RespData struct {
		NamespaceUsers model.NamespaceUsers `json:"list"`
	}
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	namespaceUsers, err := model.NamespaceUser{NamespaceID: id}.GetBindUserListByNamespaceID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{Data: RespData{NamespaceUsers: namespaceUsers}}
}

// Add one Namespace
func (namespace Namespace) Add(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	id, err := model.Namespace{Name: reqData.Name}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	err = model.NamespaceUser{NamespaceID: id}.AddAdminByNamespaceID()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{Data: RespData{ID: id}}
}

// Edit one Namespace
func (namespace Namespace) Edit(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"gt=0"`
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.Namespace{
		ID:   reqData.ID,
		Name: reqData.Name,
	}.EditRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}

// AddUser one project
func (namespace Namespace) AddUser(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		NamespaceID int64   `json:"namespaceId" validate:"gt=0"`
		UserIDs     []int64 `json:"userIds" validate:"required"`
		Role        string  `json:"role" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	namespaceUsersModel := model.NamespaceUsers{}
	for _, userID := range reqData.UserIDs {
		namespaceUserModel := model.NamespaceUser{
			NamespaceID: reqData.NamespaceID,
			UserID:      userID,
			Role:        reqData.Role,
		}
		namespaceUsersModel = append(namespaceUsersModel, namespaceUserModel)
	}

	if err := namespaceUsersModel.AddMany(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if reqData.Role == core.RoleManager {
		err := model.ProjectUser{}.AddNamespaceProjectInUserID(reqData.NamespaceID, reqData.UserIDs)
		if err != nil {
			return &core.Response{Code: core.Error, Message: err.Error()}
		}
	}

	return &core.Response{}
}

// RemoveUser one Project
func (namespace Namespace) RemoveUser(_ http.ResponseWriter, gp *core.Goploy) *core.Response {
	type ReqData struct {
		NamespaceUserID int64 `json:"namespaceUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	err := model.NamespaceUser{
		ID: reqData.NamespaceUserID,
	}.DeleteRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}
