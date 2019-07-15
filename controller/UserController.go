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
	userData, err := model.User{Account: reqData.Account}.GetDataByAccount()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if err := userData.Vaildate(reqData.Password); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	token, err := user.createToken(userData)
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
	userData, err := model.User{ID: core.GolbalUserID}.GetData()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if userData.State != 1 {
		response := core.Response{Code: 1, Message: "账号被封停"}
		response.Json(w)
		return
	}

	data := RepData{}
	data.UserInfo.ID = core.GolbalUserID
	data.UserInfo.Name = userData.Name
	data.UserInfo.Account = userData.Account

	role, err := model.Role{ID: userData.RoleID}.GetData()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	permissions, err := model.Permission{}.GetAllByPermissionList(role.PermissionList)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	var tempPermissions model.Permissions
	for _, permission := range permissions {
		data.PermissionURI = append(data.PermissionURI, permission.URI)
		if permission.PID == 0 {
			for _, pmChild := range permissions {
				if pmChild.PID == permission.ID {
					permission.Children = append(permission.Children, pmChild)
				}
			}
			tempPermissions = append(tempPermissions, permission)
		}
	}

	data.Permission = tempPermissions

	response := core.Response{Data: data}
	response.Json(w)
}

// GetList user list
func (user *User) GetList(w http.ResponseWriter, r *http.Request) {
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
	users, err := userModel.GetList(pagination)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{User: users, Pagination: *pagination}}
	response.Json(w)
}

// GetOption user list
func (user *User) GetOption(w http.ResponseWriter, r *http.Request) {
	type RepData struct {
		User model.Users `json:"userList"`
	}
	users, err := model.User{}.GetAll()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Data: RepData{User: users}}
	response.Json(w)
}

// Add one user
func (user *User) Add(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Mobile   string `json:"mobile"`
		RoleID   uint32 `json:"roleId"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	_, err = model.User{
		Account:    reqData.Account,
		Password:   reqData.Password,
		Name:       reqData.Name,
		Mobile:     reqData.Mobile,
		RoleID:     reqData.RoleID,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}.AddRow()

	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	response := core.Response{Message: "添加成功"}
	response.Json(w)
}

// Edit one user
func (user *User) Edit(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		ID       uint32 `json:"id"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Mobile   string `json:"mobile"`
		RoleID   uint32 `json:"roleId"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	err = model.User{
		ID:         reqData.ID,
		Password:   reqData.Password,
		Name:       reqData.Name,
		Mobile:     reqData.Mobile,
		RoleID:     reqData.RoleID,
		UpdateTime: time.Now().Unix(),
	}.EditRow()

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
		OldPassword string `json:"oldPwd"`
		NewPassword string `json:"newPwd"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	userData, err := model.User{ID: core.GolbalUserID}.GetData()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if err := userData.Vaildate(reqData.OldPassword); err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}

	if err := (model.User{ID: core.GolbalUserID, Password: reqData.NewPassword}.UpdatePassword()); err != nil {
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
