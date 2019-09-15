package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"goploy/core"
	"goploy/model"

	cache "github.com/patrickmn/go-cache"
)

// User 用户字段
type User Controller

// Login user login api
func (user User) Login(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if userData.State == model.Disable {
		response := core.Response{Code: core.AccountDisabled, Message: "账号已被停用"}
		response.JSON(w)
		return
	}

	if err := userData.Vaildate(reqData.Password); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	token, err := userData.CreateToken()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	model.User{ID: userData.ID, LastLoginTime: time.Now().Unix()}.UpdateLastLoginTime()

	core.Cache.Set("userInfo:"+strconv.Itoa(int(userData.ID)), &userData, cache.DefaultExpiration)
	cookie := http.Cookie{Name: core.LoginCookieName, Value: token, Path: "/", MaxAge: 86400, HttpOnly: true}
	http.SetCookie(w, &cookie)
	response := core.Response{}
	response.JSON(w)
}

// Info get user info api
func (user User) Info(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		UserInfo struct {
			ID      int64  `json:"id"`
			Account string `json:"account"`
			Name    string `json:"name"`
			Role    string `json:"role"`
		} `json:"userInfo"`
	}
	userData, err := core.GetUserInfo(gp.TokenInfo.ID)

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	data := RespData{}
	data.UserInfo.ID = gp.TokenInfo.ID
	data.UserInfo.Name = userData.Name
	data.UserInfo.Account = userData.Account
	data.UserInfo.Role = userData.Role
	response := core.Response{Data: data}
	response.JSON(w)
}

// GetList user list
func (user User) GetList(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		User       model.Users      `json:"userList"`
		Pagination model.Pagination `json:"pagination"`
	}
	pagination, err := model.PaginationFrom(gp.URLQuery)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	users, pagination, err := model.Users{}.GetList(pagination)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{User: users, Pagination: pagination}}
	response.JSON(w)
}

// GetOption user list
func (user User) GetOption(w http.ResponseWriter, gp *core.Goploy) {
	type RespData struct {
		User model.Users `json:"userList"`
	}
	users, err := model.User{}.GetAll()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Data: RespData{User: users}}
	response.JSON(w)
}

// Add one user
func (user User) Add(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		Account        string `json:"account"`
		Password       string `json:"password"`
		Name           string `json:"name"`
		Mobile         string `json:"mobile"`
		Role           string `json:"role"`
		ManageGroupStr string `json:"manageGroupStr"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	userInfo, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil && err != sql.ErrNoRows {
		response := core.Response{Message: err.Error()}
		response.JSON(w)
	} else if userInfo != (model.User{}) {
		response := core.Response{Message: "账号已存在"}
		response.JSON(w)
	}
	_, err = model.User{
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
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Edit one user
func (user User) Edit(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID             int64  `json:"id"`
		Password       string `json:"password"`
		Name           string `json:"name"`
		Mobile         string `json:"mobile"`
		Role           string `json:"role"`
		ManageGroupStr string `json:"manageGroupStr"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	err = model.User{
		ID:             reqData.ID,
		Password:       reqData.Password,
		Name:           reqData.Name,
		Mobile:         reqData.Mobile,
		Role:           reqData.Role,
		ManageGroupStr: reqData.ManageGroupStr,
		UpdateTime:     time.Now().Unix(),
	}.EditRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	response := core.Response{Message: "添加成功"}
	response.JSON(w)
}

// Remove one User
func (user User) Remove(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		ID int64 `json:"id"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if reqData.ID == 1 {
		response := core.Response{Code: core.Deny, Message: "请勿删除超管账号"}
		response.JSON(w)
		return
	}

	err = model.User{
		ID:         reqData.ID,
		UpdateTime: time.Now().Unix(),
	}.RemoveRow()

	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "删除成功"}
	response.JSON(w)
}

// ChangePassword doc
func (user User) ChangePassword(w http.ResponseWriter, gp *core.Goploy) {
	type ReqData struct {
		OldPassword string `json:"oldPwd"`
		NewPassword string `json:"newPwd"`
	}
	var reqData ReqData
	err := json.Unmarshal(gp.Body, &reqData)
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	userData, err := model.User{ID: gp.TokenInfo.ID}.GetData()
	if err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if err := userData.Vaildate(reqData.OldPassword); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}

	if err := (model.User{ID: gp.TokenInfo.ID, Password: reqData.NewPassword}.UpdatePassword()); err != nil {
		response := core.Response{Code: core.Deny, Message: err.Error()}
		response.JSON(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.JSON(w)
}
