package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"goploy/core"
	"goploy/model"

	"github.com/patrickmn/go-cache"
)

// User 用户字段
type User Controller

// Login user login api
func (user User) Login(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type ReqData struct {
		Account  string `json:"account" validate:"min=5,max=12"`
		Password string `json:"password" validate:"password"`
	}
	type RespData struct {
		Token string `json:"token"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := userData.Vaildate(reqData.Password); err != nil {
		return core.Response{Code: core.Deny, Message: err.Error()}
	}

	if userData.State == model.Disable {
		return core.Response{Code: core.AccountDisabled, Message: "Account is disabled"}
	}

	token, err := userData.CreateToken()
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	model.User{ID: userData.ID, LastLoginTime: time.Now().Unix()}.UpdateLastLoginTime()

	core.Cache.Set("userInfo:"+strconv.Itoa(int(userData.ID)), &userData, cache.DefaultExpiration)
	cookie := http.Cookie{Name: core.LoginCookieName, Value: token, Path: "/", MaxAge: 86400, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return core.Response{Data: RespData{Token: token}}
}

// Info get user info api
func (user User) Info(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type RespData struct {
		UserInfo struct {
			ID      int64  `json:"id"`
			Account string `json:"account"`
			Name    string `json:"name"`
			Role    string `json:"role"`
		} `json:"userInfo"`
	}
	data := RespData{}
	data.UserInfo.ID = gp.UserInfo.ID
	data.UserInfo.Name = gp.UserInfo.Name
	data.UserInfo.Account = gp.UserInfo.Account
	data.UserInfo.Role = gp.UserInfo.Role
	return core.Response{Data: data}
}

// GetList user list
func (user User) GetList(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type RespData struct {
		User       model.Users      `json:"userList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	users, pagination, err := model.Users{}.GetList(pagination)
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{Data: RespData{User: users, Pagination: pagination}}
}

// GetOption user list
func (user User) GetOption(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type RespData struct {
		User model.Users `json:"userList"`
	}
	users, err := model.User{}.GetAll()
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{Data: RespData{User: users}}
}

// GetProjectOption user list
func (user User) GetCanBindProjectUser(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type RespData struct {
		User model.Users `json:"userList"`
	}
	users, err := model.User{}.GetCanBindProjectUser()
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{Data: RespData{User: users}}
}

// Add one user
func (user User) Add(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type ReqData struct {
		Account        string  `json:"account" validate:"min=5,max=12"`
		Password       string  `json:"password" validate:"omitempty,password"`
		Name           string  `json:"name" validate:"required"`
		Mobile         string  `json:"mobile" validate:"omitempty,len=11,numeric"`
		Role           string  `json:"role" validate:"role"`
		ManageGroupStr string  `json:"manageGroupStr"`
		ProjectIDs     []int64 `json:"projectIds"`
	}
	type RespData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	userInfo, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil && err != sql.ErrNoRows {
		return core.Response{Code: core.Error, Message: err.Error()}
	} else if userInfo != (model.User{}) {
		return core.Response{Code: core.Error, Message: "Account is already exist"}
	}
	id, err := model.User{
		Account:        reqData.Account,
		Password:       reqData.Password,
		Name:           reqData.Name,
		Mobile:         reqData.Mobile,
		Role:           reqData.Role,
		ManageGroupStr: reqData.ManageGroupStr,
		CreateTime:     time.Now().Unix(),
		UpdateTime:     time.Now().Unix(),
	}.AddRow()

	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	projectUsersModel := model.ProjectUsers{}
	for _, projectID := range reqData.ProjectIDs {
		projectUserModel := model.ProjectUser{
			ProjectID:  projectID,
			UserID:     id,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}
	if err := projectUsersModel.AddMany(); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{Data: RespData{ID: id}}
}

// Edit one user
func (user User) Edit(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type ReqData struct {
		ID             int64   `json:"id" validate:"gt=0"`
		Password       string  `json:"password" validate:"omitempty,password"`
		Name           string  `json:"name" validate:"required"`
		Mobile         string  `json:"mobile" validate:"omitempty,len=11,numeric"`
		Role           string  `json:"role" validate:"role"`
		ManageGroupStr string  `json:"manageGroupStr"`
		ProjectIDs     []int64 `json:"projectIds"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	err := model.User{
		ID:             reqData.ID,
		Password:       reqData.Password,
		Name:           reqData.Name,
		Mobile:         reqData.Mobile,
		Role:           reqData.Role,
		ManageGroupStr: reqData.ManageGroupStr,
		UpdateTime:     time.Now().Unix(),
	}.EditRow()

	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	err = model.ProjectUser{
		UserID: reqData.ID,
	}.DeleteByUserID()

	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	projectUsersModel := model.ProjectUsers{}
	for _, projectID := range reqData.ProjectIDs {
		projectUserModel := model.ProjectUser{
			ProjectID:  projectID,
			UserID:     reqData.ID,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		projectUsersModel = append(projectUsersModel, projectUserModel)
	}
	if err := projectUsersModel.AddMany(); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{}
}

// Remove one User
func (user User) Remove(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type ReqData struct {
		ID int64 `json:"id" validate:"gt=0"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	if reqData.ID == 1 {
		return core.Response{Code: core.Error, Message: "Can not delete the super manager"}
	}

	err := model.User{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.RemoveRow()

	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{}
}

// ChangePassword doc
func (user User) ChangePassword(w http.ResponseWriter, gp *core.Goploy) core.Response {
	type ReqData struct {
		OldPassword string `json:"oldPwd" validate:"password"`
		NewPassword string `json:"newPwd" validate:"password"`
	}
	var reqData ReqData
	if err := verify(gp.Body, &reqData); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	userData, err := model.User{ID: gp.UserInfo.ID}.GetData()
	if err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := userData.Vaildate(reqData.OldPassword); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}

	if err := (model.User{ID: gp.UserInfo.ID, Password: reqData.NewPassword}.UpdatePassword()); err != nil {
		return core.Response{Code: core.Error, Message: err.Error()}
	}
	return core.Response{}
}
