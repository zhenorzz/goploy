package controller

import (
	"database/sql"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/response"
	"net/http"
	"time"

	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

type User Controller

func (u User) Routes() []core.Route {
	return []core.Route{
		core.NewWhiteRoute("/user/login", http.MethodPost, u.Login).LogFunc(middleware.AddLoginLog),
		core.NewRoute("/user/info", http.MethodGet, u.Info),
		core.NewRoute("/user/getList", http.MethodGet, u.GetList),
		core.NewRoute("/user/getTotal", http.MethodGet, u.GetTotal),
		core.NewRoute("/user/add", http.MethodPost, u.Add).Roles(core.RoleAdmin),
		core.NewRoute("/user/edit", http.MethodPut, u.Edit).Roles(core.RoleAdmin),
		core.NewRoute("/user/remove", http.MethodDelete, u.Remove).Roles(core.RoleAdmin),
		core.NewRoute("/user/changePassword", http.MethodPut, u.ChangePassword),
	}
}

func (User) Login(gp *core.Goploy) core.Response {
	type ReqData struct {
		Account  string `json:"account" validate:"min=5,max=25"`
		Password string `json:"password" validate:"password"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if config.Toml.LDAP.Enabled && userData.ID != 1 {
		conn, err := ldap.DialURL(config.Toml.LDAP.URL)
		if err != nil {
			return response.JSON{Code: response.Deny, Message: err.Error()}
		}

		if config.Toml.LDAP.BindDN != "" {
			if err := conn.Bind(config.Toml.LDAP.BindDN, config.Toml.LDAP.Password); err != nil {
				return response.JSON{Code: response.Deny, Message: err.Error()}
			}
		}

		filter := fmt.Sprintf("(%s=%s)", config.Toml.LDAP.UID, reqData.Account)
		if config.Toml.LDAP.UserFilter != "" {
			filter = fmt.Sprintf("(&(%s)%s)", config.Toml.LDAP.UserFilter, filter)
		}

		searchRequest := ldap.NewSearchRequest(
			config.Toml.LDAP.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			filter,
			[]string{config.Toml.LDAP.UID},
			nil)

		sr, err := conn.Search(searchRequest)
		if err != nil {
			return response.JSON{Code: response.Deny, Message: err.Error()}
		}
		if len(sr.Entries) != 1 {
			return response.JSON{Code: response.Deny, Message: fmt.Sprintf("No %s record in baseDN %s", reqData.Account, config.Toml.LDAP.BaseDN)}
		}
		if err := conn.Bind(sr.Entries[0].DN, reqData.Password); err != nil {
			return response.JSON{Code: response.Deny, Message: err.Error()}
		}
	} else {
		if err := userData.Validate(reqData.Password); err != nil {
			return response.JSON{Code: response.Deny, Message: err.Error()}
		}
	}

	if userData.State == model.Disable {
		return response.JSON{Code: response.AccountDisabled, Message: "Account is disabled"}
	}
	namespaceList, err := model.Namespace{UserID: userData.ID}.GetAllByUserID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	} else if len(namespaceList) == 0 {
		return response.JSON{Code: response.Error, Message: "No space assigned, please contact the administrator"}
	}

	token, err := userData.CreateToken()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	_ = model.User{ID: userData.ID, LastLoginTime: time.Now().Format("20060102150405")}.UpdateLastLoginTime()

	cookie := http.Cookie{
		Name:     config.Toml.Cookie.Name,
		Value:    token,
		Path:     "/",
		MaxAge:   config.Toml.Cookie.Expire,
		HttpOnly: true,
	}
	http.SetCookie(gp.ResponseWriter, &cookie)
	return response.JSON{
		Data: struct {
			Token         string           `json:"token"`
			NamespaceList model.Namespaces `json:"namespaceList"`
		}{Token: token, NamespaceList: namespaceList},
	}
}

func (User) Info(gp *core.Goploy) core.Response {
	type RespData struct {
		UserInfo struct {
			ID           int64  `json:"id"`
			Account      string `json:"account"`
			Name         string `json:"name"`
			SuperManager int64  `json:"superManager"`
		} `json:"userInfo"`
	}
	data := RespData{}
	data.UserInfo.ID = gp.UserInfo.ID
	data.UserInfo.Name = gp.UserInfo.Name
	data.UserInfo.Account = gp.UserInfo.Account
	data.UserInfo.SuperManager = gp.UserInfo.SuperManager
	return response.JSON{Data: data}
}

func (User) GetList(gp *core.Goploy) core.Response {
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	users, err := model.User{}.GetList(pagination)
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Users model.Users `json:"list"`
		}{Users: users},
	}
}

func (User) GetTotal(*core.Goploy) core.Response {
	total, err := model.User{}.GetTotal()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Total int64 `json:"total"`
		}{Total: total},
	}
}

func (User) Add(gp *core.Goploy) core.Response {
	type ReqData struct {
		Account      string `json:"account" validate:"min=5,max=25"`
		Password     string `json:"password" validate:"omitempty,password"`
		Name         string `json:"name" validate:"required"`
		Contact      string `json:"contact" validate:"omitempty,len=11,numeric"`
		SuperManager int64  `json:"superManager" validate:"min=0,max=1"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	userInfo, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil && err != sql.ErrNoRows {
		return response.JSON{Code: response.Error, Message: err.Error()}
	} else if userInfo != (model.User{}) {
		return response.JSON{Code: response.Error, Message: "Account is already exist"}
	}
	id, err := model.User{
		Account:      reqData.Account,
		Password:     reqData.Password,
		Name:         reqData.Name,
		Contact:      reqData.Contact,
		SuperManager: reqData.SuperManager,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.SuperManager == model.SuperManager {
		if err := (model.NamespaceUser{UserID: id}).AddAdminByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		if err := (model.ProjectUser{UserID: id}).AddAdminByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: id},
	}
}

func (User) Edit(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"gt=0"`
		Password     string `json:"password" validate:"omitempty,password"`
		Name         string `json:"name" validate:"required"`
		Contact      string `json:"contact" validate:"omitempty,len=11,numeric"`
		SuperManager int64  `json:"superManager" validate:"min=0,max=1"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	userInfo, err := model.User{ID: reqData.ID}.GetData()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err = model.User{
		ID:           reqData.ID,
		Password:     reqData.Password,
		Name:         reqData.Name,
		Contact:      reqData.Contact,
		SuperManager: reqData.SuperManager,
	}.EditRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if userInfo.SuperManager == model.SuperManager && reqData.SuperManager == model.GeneralUser {
		if err := (model.NamespaceUser{UserID: reqData.ID}).DeleteByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		if err := (model.ProjectUser{UserID: reqData.ID}).DeleteByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	} else if userInfo.SuperManager == model.GeneralUser && reqData.SuperManager == model.SuperManager {
		if err := (model.NamespaceUser{UserID: reqData.ID}).AddAdminByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
		if err := (model.ProjectUser{UserID: reqData.ID}).AddAdminByUserID(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	return response.JSON{}
}

func (User) Remove(gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	if reqData.ID == 1 {
		return response.JSON{Code: response.Error, Message: "Can not delete the super manager"}
	}
	if err := (model.User{ID: reqData.ID}).RemoveRow(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (User) ChangePassword(gp *core.Goploy) core.Response {
	type ReqData struct {
		OldPassword string `json:"oldPwd" validate:"password"`
		NewPassword string `json:"newPwd" validate:"password"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	userData, err := model.User{ID: gp.UserInfo.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := userData.Validate(reqData.OldPassword); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := (model.User{ID: gp.UserInfo.ID, Password: reqData.NewPassword}).UpdatePassword(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
