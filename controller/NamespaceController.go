package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/permission"
	"github.com/zhenorzz/goploy/response"
	"net/http"
	"strconv"
)

type Namespace Controller

func (n Namespace) Routes() []core.Route {
	return []core.Route{
		core.NewRoute("/namespace/getList", http.MethodGet, n.GetList).Permissions(permission.ShowNamespacePage),
		core.NewRoute("/namespace/getOption", http.MethodGet, n.GetOption),
		core.NewRoute("/namespace/getBindUserList", http.MethodGet, n.GetBindUserList).Permissions(permission.ShowNamespacePage),
		core.NewRoute("/namespace/getUserOption", http.MethodGet, n.GetUserOption),
		core.NewRoute("/namespace/add", http.MethodPost, n.Add).Roles(core.RoleAdmin).Permissions(permission.AddNamespace),
		core.NewRoute("/namespace/edit", http.MethodPut, n.Edit).Roles(core.RoleAdmin, core.RoleManager).Permissions(permission.EditNamespace),
		core.NewRoute("/namespace/addUser", http.MethodPost, n.AddUser).Roles(core.RoleAdmin, core.RoleManager).Permissions(permission.AddNamespaceUser),
		core.NewRoute("/namespace/removeUser", http.MethodDelete, n.RemoveUser).Roles(core.RoleAdmin, core.RoleManager).Permissions(permission.DeleteNamespaceUser),
	}
}

func (Namespace) GetList(gp *core.Goploy) core.Response {
	namespaceList, err := model.Namespace{UserID: gp.UserInfo.ID}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			Namespaces model.Namespaces `json:"list"`
		}{Namespaces: namespaceList},
	}
}

func (Namespace) GetOption(gp *core.Goploy) core.Response {
	namespaceUsers, err := model.NamespaceUser{UserID: gp.UserInfo.ID}.GetUserNamespaceList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) GetUserOption(gp *core.Goploy) core.Response {
	namespaceUsers, err := model.NamespaceUser{NamespaceID: gp.Namespace.ID}.GetAllUserByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) GetBindUserList(gp *core.Goploy) core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	namespaceUsers, err := model.NamespaceUser{NamespaceID: id}.GetBindUserListByNamespaceID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			NamespaceUsers model.NamespaceUsers `json:"list"`
		}{NamespaceUsers: namespaceUsers},
	}
}

func (Namespace) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	id, err := model.Namespace{Name: reqData.Name}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{NamespaceID: id}).AddAdminByNamespaceID(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (Namespace) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID   int64  `json:"id" validate:"gt=0"`
		Name string `json:"name" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	err := model.Namespace{
		ID:   reqData.ID,
		Name: reqData.Name,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Namespace) AddUser(gp *core.Goploy) core.Response {
	type ReqData struct {
		NamespaceID int64   `json:"namespaceId" validate:"gt=0"`
		UserIDs     []int64 `json:"userIds" validate:"required"`
		RoleID      int64   `json:"roleId" validate:"required"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	namespaceUsersModel := model.NamespaceUsers{}
	for _, userID := range reqData.UserIDs {
		namespaceUserModel := model.NamespaceUser{
			NamespaceID: reqData.NamespaceID,
			UserID:      userID,
			RoleID:      reqData.RoleID,
		}
		namespaceUsersModel = append(namespaceUsersModel, namespaceUserModel)
	}

	if err := namespaceUsersModel.AddMany(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{}
}

func (Namespace) RemoveUser(gp *core.Goploy) core.Response {
	type ReqData struct {
		NamespaceUserID int64 `json:"namespaceUserId" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.NamespaceUser{ID: reqData.NamespaceUserID}).DeleteRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
