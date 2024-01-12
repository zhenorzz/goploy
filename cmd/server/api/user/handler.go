// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/wenlng/go-captcha/captcha"
	"github.com/zhenorzz/goploy/cmd/server/api"
	"github.com/zhenorzz/goploy/cmd/server/api/middleware"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/cache"
	"github.com/zhenorzz/goploy/internal/media"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/internal/validator"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User api.API

func (u User) Handler() []server.Route {
	return []server.Route{
		server.NewWhiteRoute("/user/login", http.MethodPost, u.Login).LogFunc(middleware.AddLoginLog),
		server.NewWhiteRoute("/user/extLogin", http.MethodPost, u.ExtLogin).LogFunc(middleware.AddLoginLog),
		server.NewRoute("/user/getList", http.MethodGet, u.GetList).Permissions(config.ShowMemberPage),
		server.NewRoute("/user/info", http.MethodGet, u.Info),
		server.NewRoute("/user/changePassword", http.MethodPut, u.ChangePassword),
		server.NewRoute("/user/getApiKey", http.MethodGet, u.GetApiKey),
		server.NewRoute("/user/generateApiKey", http.MethodPut, u.GenerateApiKey),
		server.NewRoute("/user/add", http.MethodPost, u.Add).Permissions(config.AddMember).LogFunc(middleware.AddOPLog),
		server.NewRoute("/user/edit", http.MethodPut, u.Edit).Permissions(config.EditMember).LogFunc(middleware.AddOPLog),
		server.NewRoute("/user/remove", http.MethodDelete, u.Remove).Permissions(config.DeleteMember).LogFunc(middleware.AddOPLog),
		server.NewWhiteRoute("/user/mediaLogin", http.MethodPost, u.MediaLogin).LogFunc(middleware.AddLoginLog),
		server.NewWhiteRoute("/user/getConfig", http.MethodGet, u.GetConfig),
		server.NewWhiteRoute("/user/getCaptcha", http.MethodGet, u.GetCaptcha),
		server.NewWhiteRoute("/user/checkCaptcha", http.MethodPost, u.CheckCaptcha),
	}
}

// Login user
// @Summary Login
// @Tags User
// @Produce json
// @Param request body user.Login.ReqData true "body params"
// @Success 200 {object} response.JSON{data=user.Login.RespData}
// @Router /user/login [post]
func (User) Login(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account     string `json:"account" validate:"required,min=1,max=25"`
		Password    string `json:"password" validate:"required,password"`
		NewPassword string `json:"newPassword"`
		CaptchaKey  string `json:"captchaKey" validate:"omitempty"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.IllegalParam, Message: err.Error()}
	}

	userCache := cache.GetUserCache()

	if config.Toml.Captcha.Enabled && userCache.IsShowCaptcha(reqData.Account) {
		captchaCache := cache.GetCaptchaCache()
		if !captchaCache.IsChecked(reqData.CaptchaKey) {
			return response.JSON{Code: response.Error, Message: "Captcha error, please check captcha again"}
		}
		// captcha should be deleted after check
		captchaCache.Delete(reqData.CaptchaKey)
	}

	if userCache.IsLock(reqData.Account) {
		return response.JSON{Code: response.Error, Message: "Your account has been locked, please retry login in 15 minutes"}
	}

	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
			[]string{config.Toml.LDAP.UID, config.Toml.LDAP.Name},
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

		// add user when user could not be found
		if userData.ID == 0 {
			userData = model.User{
				Account:  reqData.Account,
				Password: reqData.Password,
				Name:     sr.Entries[0].GetAttributeValue(config.Toml.LDAP.Name),
				State:    model.Enable,
			}
			userData.ID, err = userData.AddRow()
			return response.JSON{Code: response.Error, Message: err.Error()}
		}

	} else {
		if userData.ID == 0 {
			return response.JSON{Code: response.Error, Message: "We couldn't verify your identity. Please confirm if your username and password are correct"}
		}
		if err := userData.Validate(reqData.Password); err != nil {
			// error times over 5 times, then lock the account 15 minutes
			if userCache.IncrErrorTimes(reqData.Account, cache.UserCacheExpireTime) >= cache.UserCacheMaxErrorTimes {
				userCache.LockAccount(reqData.Account, cache.UserCacheLockTime)
			}
			return response.JSON{Code: response.Deny, Message: err.Error()}
		}
	}

	userCache.DeleteErrorTimes(reqData.Account)

	if userData.State == model.Disable {
		return response.JSON{Code: response.AccountDisabled, Message: "Account is disabled"}
	}

	namespaceList, err := model.Namespace{UserID: userData.ID}.GetAllByUserID()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	} else if len(namespaceList) == 0 {
		return response.JSON{Code: response.Error, Message: "No space assigned, please contact the administrator"}
	}

	if reqData.NewPassword != "" {
		if err := validator.Validate.Var(reqData.NewPassword, "password"); err != nil {
			return response.JSON{Code: response.PasswordExpired, Message: err.Error()}
		}
		if reqData.Password == reqData.NewPassword {
			return response.JSON{Code: response.PasswordExpired, Message: "The password cannot be the same as the previous one"}
		}
		if err := (model.User{ID: userData.ID, Password: reqData.NewPassword, PasswordUpdateTime: sql.NullString{
			String: time.Now().Format("20060102150405"),
			Valid:  true,
		}}).UpdatePassword(); err != nil {
			return response.JSON{Code: response.Error, Message: err.Error()}
		}
	} else if !userData.PasswordUpdateTime.Valid {
		return response.JSON{Code: response.PasswordExpired, Message: "You need to change your password upon first login"}
	} else if config.Toml.APP.PasswordPeriod > 0 {
		passwordUpdateTime, _ := time.Parse(time.DateTime, userData.PasswordUpdateTime.String)
		passwordUpdateTime = passwordUpdateTime.AddDate(0, 0, config.Toml.APP.PasswordPeriod)
		if passwordUpdateTime.Before(time.Now()) {
			return response.JSON{Code: response.PasswordExpired, Message: "Password expired, please change"}
		}
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
	gp.UserInfo = userData // for add login log
	type RespData struct {
		Token         string           `json:"token"`
		NamespaceList model.Namespaces `json:"namespaceList"`
	}

	return response.JSON{
		Data: RespData{Token: token, NamespaceList: namespaceList},
	}
}

// ExtLogin user
// @Summary External login
// @Tags User
// @Produce json
// @Param request body user.ExtLogin.ReqData true "body params"
// @Success 200 {object} response.JSON{data=user.ExtLogin.RespData}
// @Router /user/extLogin [post]
func (User) ExtLogin(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account string `json:"account" validate:"required,min=1,max=25"`
		Time    int64  `json:"time"`
		Token   string `json:"token"  validate:"required,len=32"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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
	gp.UserInfo = userData // for add login log
	type RespData struct {
		Token         string           `json:"token"`
		NamespaceList model.Namespaces `json:"namespaceList"`
	}
	return response.JSON{
		Data: RespData{Token: token, NamespaceList: namespaceList},
	}
}

// Info shows user information
// @Summary Show logged-in user information
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=user.Info.RespData}
// @Router /user/info [get]
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

// GetList lists all users
// @Summary List all users
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Success 200 {object} response.JSON{data=user.GetList.RespData}
// @Router /user/getList [get]
func (User) GetList(*server.Goploy) server.Response {
	users, err := model.User{}.GetList()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	type RespData struct {
		Users model.Users `json:"list"`
	}
	return response.JSON{
		Data: RespData{Users: users},
	}
}

// Add adds a user
// @Summary Add a user
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body user.Add.ReqData true "body params"
// @Success 200 {object} response.JSON{data=user.Add.RespData}
// @Router /user/add [post]
func (User) Add(gp *server.Goploy) server.Response {
	type ReqData struct {
		Account      string `json:"account" validate:"required,min=1,max=25"`
		Password     string `json:"password" validate:"omitempty,password"`
		Name         string `json:"name" validate:"required"`
		Contact      string `json:"contact" validate:"omitempty,len=11,numeric"`
		SuperManager int64  `json:"superManager" validate:"oneof=0 1"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	userInfo, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

	type RespData struct {
		ID int64 `json:"id"`
	}
	return response.JSON{
		Data: RespData{ID: id},
	}
}

// Edit edits the user
// @Summary Edit the user
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body user.Edit.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /user/edit [put]
func (User) Edit(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID           int64  `json:"id" validate:"required,gt=0"`
		Password     string `json:"password" validate:"omitempty,password"`
		Name         string `json:"name" validate:"required"`
		Contact      string `json:"contact" validate:"omitempty,len=11,numeric"`
		SuperManager int64  `json:"superManager" validate:"oneof=0 1"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// Remove removes the user
// @Summary Remove the user
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body user.Remove.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /user/remove [delete]
func (User) Remove(gp *server.Goploy) server.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"required,gt=0"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
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

// ChangePassword change the user password
// @Summary Change the user password
// @Tags User
// @Produce json
// @Security ApiKeyHeader || ApiKeyQueryParam || NamespaceHeader || NamespaceQueryParam
// @Param request body user.ChangePassword.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /user/changePassword [put]
func (User) ChangePassword(gp *server.Goploy) server.Response {
	type ReqData struct {
		OldPassword string `json:"oldPwd" validate:"required,password"`
		NewPassword string `json:"newPwd" validate:"required,password"`
	}
	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	userData, err := model.User{ID: gp.UserInfo.ID}.GetData()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if err := userData.Validate(reqData.OldPassword); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	if reqData.OldPassword == reqData.NewPassword {
		return response.JSON{Code: response.Error, Message: "The password cannot be the same as the previous one."}
	}

	if err := (model.User{ID: gp.UserInfo.ID, Password: reqData.NewPassword, PasswordUpdateTime: sql.NullString{
		String: time.Now().Format("20060102150405"),
		Valid:  true,
	}}).UpdatePassword(); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (User) GetApiKey(gp *server.Goploy) server.Response {
	apiKey, err := model.User{ID: gp.UserInfo.ID}.GetApiKey()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{Data: apiKey}
}

func (User) GenerateApiKey(gp *server.Goploy) server.Response {
	apiKey, err := model.User{ID: gp.UserInfo.ID}.CreateApiKey()
	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	return response.JSON{Data: apiKey}
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
	if errors.Is(err, sql.ErrNoRows) {
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
	gp.UserInfo = userData // for add login log
	return response.JSON{
		Data: struct {
			Token         string           `json:"token"`
			NamespaceList model.Namespaces `json:"namespaceList"`
		}{Token: token, NamespaceList: namespaceList},
	}
}

// GetConfig show the captcha config
// @Summary Show the captcha config
// @Tags User
// @Produce json
// @Success 200 {object} response.JSON{data=user.GetConfig.RespData}
// @Router /user/getConfig [get]
func (User) GetConfig(*server.Goploy) server.Response {

	dingtalk, feishu := "", ""

	if config.Toml.Dingtalk.AppKey != "" && config.Toml.Dingtalk.AppSecret != "" {
		dingtalk = fmt.Sprintf(
			"https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=%s&scope=openid&prompt=consent&state=dingtalk",
			config.Toml.Dingtalk.AppKey,
		)
	}

	if config.Toml.Feishu.AppKey != "" && config.Toml.Feishu.AppSecret != "" {
		feishu = fmt.Sprintf(
			"https://passport.feishu.cn/suite/passport/oauth/authorize?client_id=%s&response_type=code&state=feishu",
			config.Toml.Feishu.AppKey,
		)
	}

	type RespData struct {
		LDAP struct {
			Enabled bool `json:"enabled"`
		} `json:"ldap"`
		Captcha struct {
			Enabled bool `json:"enabled"`
		} `json:"captcha"`
		MediaURL struct {
			Dingtalk string `json:"dingtalk"`
			Feishu   string `json:"feishu"`
		} `json:"mediaURL"`
	}
	resp := &RespData{}
	resp.LDAP.Enabled = config.Toml.LDAP.Enabled
	resp.Captcha.Enabled = config.Toml.Captcha.Enabled
	resp.MediaURL.Dingtalk = dingtalk
	resp.MediaURL.Feishu = feishu

	return response.JSON{
		Data: resp,
	}
}

// GetCaptcha show a captcha
// @Summary Show a captcha
// @Tags User
// @Produce json
// @Param request query user.GetCaptcha.ReqData true "query params"
// @Success 200 {object} response.JSON{data=user.GetCaptcha.RespData}
// @Router /user/getCaptcha [get]
func (User) GetCaptcha(gp *server.Goploy) server.Response {
	type ReqData struct {
		Language string `schema:"language" validate:"required"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	capt := captcha.GetCaptcha()
	if reqData.Language == "zh-cn" {
		chars := captcha.GetCaptchaDefaultChars()
		_ = capt.SetRangChars(*chars)
	} else {
		chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		_ = capt.SetRangChars(strings.Split(chars, ""))
	}

	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		return response.JSON{Code: response.AccountDisabled, Message: "generate captcha fail, error msg:" + err.Error()}
	}

	captchaCache := cache.GetCaptchaCache()
	captchaCache.Set(key, dots, 2*time.Minute)

	type RespData struct {
		Base64      string `json:"base64"`
		ThumbBase64 string `json:"thumbBase64"`
		Key         string `json:"key"`
	}
	return response.JSON{
		Data: RespData{
			Base64:      b64,
			ThumbBase64: tb64,
			Key:         key,
		},
	}
}

// CheckCaptcha check the captcha
// @Summary Check the captcha
// @Tags User
// @Produce json
// @Param request body user.CheckCaptcha.ReqData true "body params"
// @Success 200 {object} response.JSON
// @Router /user/checkCaptcha [post]
func (User) CheckCaptcha(gp *server.Goploy) server.Response {
	type ReqData struct {
		Key         string  `json:"key" validate:"required"`
		Dots        []int64 `json:"dots" validate:"required"`
		RedirectUri string  `json:"redirectUri" validate:"omitempty"`
	}

	var reqData ReqData
	if err := gp.Decode(&reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	captchaCache := cache.GetCaptchaCache()
	dotsCache, ok := captchaCache.Get(reqData.Key)
	if !ok {
		return response.JSON{Code: response.Error, Message: "Illegal key, please refresh the captcha again"}
	}
	dots, _ := dotsCache.(map[int]captcha.CharDot)

	check := false
	if (len(dots) * 2) == len(reqData.Dots) {
		for i, dot := range dots {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", reqData.Dots[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", reqData.Dots[k]), 64)

			check = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 15)
			if !check {
				break
			}
		}
	}

	if !check {
		return response.JSON{Code: response.Error, Message: "check captcha fail"}
	}

	// set captcha key checked for login verify
	captchaCache.Set(reqData.Key, true, 2*time.Minute)

	return response.JSON{}
}
