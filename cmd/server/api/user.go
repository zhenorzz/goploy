// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/internal/media"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zhenorzz/goploy/config"
)

type User API

func (u User) Handler() []server.Route {
	return []server.Route{
		server.NewWhiteRoute("/user/login", http.MethodPost, u.Login).LogFunc(middleware.AddLoginLog),
		server.NewWhiteRoute("/user/extLogin", http.MethodPost, u.ExtLogin).LogFunc(middleware.AddLoginLog),
		server.NewRoute("/user/info", http.MethodGet, u.Info),
		server.NewRoute("/user/changePassword", http.MethodPut, u.ChangePassword),
		server.NewRoute("/user/getList", http.MethodGet, u.GetList).Permissions(config.ShowMemberPage),
		server.NewRoute("/user/getOption", http.MethodGet, u.GetOption),
		server.NewRoute("/user/add", http.MethodPost, u.Add).Permissions(config.AddMember).LogFunc(middleware.AddOPLog),
		server.NewRoute("/user/edit", http.MethodPut, u.Edit).Permissions(config.EditMember).LogFunc(middleware.AddOPLog),
		server.NewRoute("/user/remove", http.MethodDelete, u.Remove).Permissions(config.DeleteMember).LogFunc(middleware.AddOPLog),
		server.NewWhiteRoute("/user/mediaLogin", http.MethodPost, u.MediaLogin),
		server.NewWhiteRoute("/user/getMediaLoginUrl", http.MethodGet, u.GetMediaLoginUrl),
	}
}

func (User) Login(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account  string `json:"account" validate:"min=1,max=25"`
		Password string `json:"password" validate:"password"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err == sql.ErrNoRows {
		return response.JSON{Code: response.Error, Message: "We couldn't verify your identity. Please confirm if your username and password are correct."}
	}
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

func (User) ExtLogin(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account string `json:"account" validate:"min=1,max=25"`
		Time    int64  `json:"time"`
		Token   string `json:"token"  validate:"len=32"`
	}
	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	if time.Now().Unix() > reqData.Time+30 {
		return response.JSON{Code: response.IllegalParam, Message: "request time expired"}
	}

	h := md5.New()
	h.Write([]byte(reqData.Account + config.Toml.JWT.Key + strconv.FormatInt(reqData.Time, 10)))
	signedToken := hex.EncodeToString(h.Sum(nil))

	if signedToken != reqData.Token {
		return response.JSON{Code: response.IllegalParam, Message: "sign error"}
	}

	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
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

func (User) Info(gp *server.Goploy) server.Response {
	type RespData struct {
		UserInfo struct {
			ID           int64  `json:"id"`
			Account      string `json:"account"`
			Name         string `json:"name"`
			SuperManager int64  `json:"superManager"`
		} `json:"userInfo"`
		Namespace struct {
			ID            int64   `json:"id"`
			PermissionIDs []int64 `json:"permissionIds"`
		} `json:"namespace"`
	}
	data := RespData{}
	data.UserInfo.ID = gp.UserInfo.ID
	data.UserInfo.Name = gp.UserInfo.Name
	data.UserInfo.Account = gp.UserInfo.Account
	data.UserInfo.SuperManager = gp.UserInfo.SuperManager

	data.Namespace.ID = gp.Namespace.ID
	data.Namespace.PermissionIDs = make([]int64, len(gp.Namespace.PermissionIDs))
	i := 0
	for k := range gp.Namespace.PermissionIDs {
		data.Namespace.PermissionIDs[i] = k
		i++
	}

	return response.JSON{Data: data}
}

func (User) GetList(*server.Goploy) server.Response {
	users, err := model.User{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			Users model.Users `json:"list"`
		}{Users: users},
	}
}

func (User) GetOption(*server.Goploy) server.Response {
	users, err := model.User{}.GetAll()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{Data: struct {
		Users model.Users `json:"list"`
	}{Users: users}}
}

func (User) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account      string `json:"account" validate:"min=1,max=25"`
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

func (User) Edit(gp *server.Goploy) server.Response {
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

func (User) Remove(gp *server.Goploy) server.Response {
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

func (User) ChangePassword(gp *server.Goploy) server.Response {
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

func (User) MediaLogin(gp *server.Goploy) server.Response {
	type ReqData struct {
		AuthCode    string `json:"authCode" validate:"required"`
		State       string `json:"state" validate:"required"`
		RedirectUri string `json:"redirectUri" validate:"omitempty"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	var mobile string
	var err error

	mediaService := media.GetMedia(reqData.State)
	if mobile, err = mediaService.Login(reqData.AuthCode, reqData.RedirectUri); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	userData, err := model.User{Contact: mobile}.GetDataByContact()
	if err == sql.ErrNoRows {
		return response.JSON{Code: response.Error, Message: "We couldn't find your account, please contact admin"}
	}
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
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

func (User) GetMediaLoginUrl(gp *server.Goploy) server.Response {
	type ReqData struct {
		RedirectUri string `schema:"redirectUri" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	reqData.RedirectUri = url.QueryEscape(reqData.RedirectUri)

	dingtalk, feishu := "", ""

	if config.Toml.Dingtalk.AppKey != "" {
		dingtalk = fmt.Sprintf(
			"https://login.dingtalk.com/oauth2/auth?redirect_uri=%s&response_type=code&client_id=%s&scope=openid&prompt=consent&state=dingtalk",
			reqData.RedirectUri,
			config.Toml.Dingtalk.AppKey,
		)
	}

	if config.Toml.Feishu.AppKey != "" {
		feishu = fmt.Sprintf(
			"https://passport.feishu.cn/suite/passport/oauth/authorize?redirect_uri=%s&client_id=%s&response_type=code&state=feishu",
			reqData.RedirectUri,
			config.Toml.Feishu.AppKey,
		)
	}

	return response.JSON{
		Data: struct {
			Dingtalk string `json:"dingtalk"`
			Feishu   string `json:"feishu"`
		}{Dingtalk: dingtalk, Feishu: feishu},
	}
}
