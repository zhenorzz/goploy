package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goploy/core"
	"goploy/model"
	"goploy/route"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestApi(t *testing.T) {
	godotenv.Load("../.env")
	model.Init()
	// user login have token
	userLogin(t)
	t.Run("user/info", userInfo)
	t.Run("user/getList", getUserList)
	t.Run("user/getOption", getUserOption)
	t.Run("user/add", addUser)
	t.Run("user/edit", editUser)
	t.Run("user/remove", removeUser)
	t.Run("user/changePassword", changeUserPassword)
}

var handler = route.Init()
var token = ""

func request(t *testing.T, method, url string, body interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if token != "" {
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", core.LoginCookieName, token))
	}
	r := httptest.NewRecorder()
	handler.ServeHTTP(r, req)
	// 检测返回的状态码
	if r.Code != http.StatusOK {
		t.Fatalf("http request error, code: %d", r.Code)
	}

	var resp core.Response

	// 检测返回的json格式
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err.Error())
	}

	// 检测接口返回值
	if resp.Code != core.Pass {
		t.Fatalf("http response error, content: %v", resp)
	}

	if url == "/user/login" {
		type RespData struct {
			Token string `json:"token"`
		}
		token = resp.Data.(map[string]interface{})["token"].(string)
	}
}
