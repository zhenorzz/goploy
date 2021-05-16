package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"strconv"
)

// Namespace struct
type Namespace Controller

// GetList namespace list
func (Namespace) GetList(gp *core.Goploy) *core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	namespaceList, err := model.Namespace{UserID: gp.UserInfo.ID}.GetListByUserID(pagination)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			Namespaces model.Namespaces `json:"list"`
		}{Namespaces: namespaceList},
	}
}

// GetTotal namespace list
func (Namespace) GetTotal(gp *core.Goploy) *core.Response {
	total, err := model.Namespace{UserID: gp.UserInfo.ID}.GetTotalByUserID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

// GetUserOption user list
func (Namespace) GetUserOption(gp *core.Goploy) *core.Response {
	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllUserByNamespaceID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

// GetBindUserList user list
func (Namespace) GetBindUserList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	namespaceUsers, err := model.NamespaceUser{NamespaceID: id}.GetBindUserListByNamespaceID()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

// Add one namespace
func (Namespace) Add(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	id, err := model.Namespace{Name: reqData.Name}.AddRow()

	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{NamespaceID: id}).AddAdminByNamespaceID(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

// Edit one namespace
func (Namespace) Edit(gp *core.Goploy) *core.Response {
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

// AddUser to namespace
func (Namespace) AddUser(gp *core.Goploy) *core.Response {
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

// RemoveUser from namespace
func (Namespace) RemoveUser(gp *core.Goploy) *core.Response {
	type ReqData struct {
		NamespaceUserID int64 `json:"namespaceUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{ID: reqData.NamespaceUserID}).DeleteRow(); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	return &core.Response{}
}
