package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// User 用户字段
type User struct{}

// IsShowPhrase is show phrase
func (user *User) IsShowPhrase(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		Show bool `json:"show"`
	}
	data := RepData{Show: false}
	response := core.Response{Data: data}
	response.Json(w)
}

// Login user login api
func (user *User) Login(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	type RepData struct {
		Token string `json:"token"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.User{Account: reqData.Account, Password: reqData.Password}
	err = model.Vaildate()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	token, err := user.createToken(model)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	data := RepData{Token: token}
	response := core.Response{Data: data}
	response.Json(w)
}

// Info get user info api
func (user *User) Info(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		UserInfo struct {
			ID      uint32 `json:"id"`
			Account string `json:"account"`
			Name    string `json:"name"`
		} `json:"userInfo"`
		Permission    model.Permissions `json:"permission"`
		PermissionURI []string          `json:"permissionUri"`
	}
	userModel := model.User{ID: core.GolbalUserID}
	err := userModel.QueryRow()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	data := RepData{}
	data.UserInfo.ID = core.GolbalUserID
	data.UserInfo.Name = userModel.Name
	data.UserInfo.Account = userModel.Account

	roleModel := model.Role{ID: userModel.RoleID}
	err = roleModel.QueryRow()

	permissionsModel := model.Permissions{}
	err = permissionsModel.Query(roleModel.PermissionList)
	var permissions model.Permissions
	for _, permission := range permissionsModel {
		data.PermissionURI = append(data.PermissionURI, permission.URI)
		if permission.PID == 0 {
			for _, pmChild := range permissionsModel {
				if pmChild.PID == permission.ID {
					permission.Children = append(permission.Children, pmChild)
				}
			}
			permissions = append(permissions, permission)
		}
	}

	data.Permission = permissions

	response := core.Response{Data: data}
	response.Json(w)
}

// Get user list
func (user *User) Get(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		User       model.Users      `json:"userList"`
		Pagination model.Pagination `json:"pagination"`
	}
	userModel := model.Users{}
	pagination, err := model.NewPagination(r.URL.Query())
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	if err := userModel.Query(pagination); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{User: userModel, Pagination: *pagination}}
	response.Json(w)
}

// Add one user
func (user *User) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Mobile   string `json:"mobile"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.User{
		Account:    reqData.Account,
		Password:   reqData.Password,
		Name:       reqData.Name,
		Mobile:     reqData.Mobile,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = model.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.Json(w)
}

// ChangePassword doc
func (user *User) ChangePassword(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Account     string `json:"account"`
		OldPassword string `json:"oldPwd"`
		NewPassword string `json:"newPwd"`
	}
	type RepData struct {
		Token string `json:"token"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	userModel := model.User{Account: reqData.Account, Password: reqData.OldPassword}
	if err := userModel.Vaildate(); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if err := userModel.UpdatePassword(reqData.NewPassword); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "修改成功"}
	response.Json(w)
}

func (user *User) createToken(u model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.ID,
		"name": u.Name,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"nbf":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))

	//Sign and get the complete encoded token as string
	return tokenString, err
}
