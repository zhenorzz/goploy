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

	"gotest.tools/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestApi(t *testing.T) {
	godotenv.Load("../.env")
	model.Init()
	// userLogin(t)
	userInfo(t)
}

var handler = route.Init()
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njg2OTIwODYsImlkIjoxLCJuYW1lIjoi6LaF566hIiwibmJmIjoxNTY4NjA1Njg2fQ.PdC3HJx5F-Ig17Oo3m2xQdO1hIzNpgCdi5BK69XLFEQ"

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
	assert.Equal(t, r.Code, http.StatusOK)
	var resp core.Response
	assert.NilError(t, json.NewDecoder(r.Body).Decode(&resp))
	assert.Equal(t, resp.Code, core.Pass, resp.Message)
	if url == "/user/login" {
		type RespData struct {
			Token string `json:"token"`
		}
		token = resp.Data.(map[string]interface{})["token"].(string)
	}
}

func userLogin(t *testing.T) {
	//创建一个请求
	body := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}{
		Account:  "admin",
		Password: "admin!@#",
	}
	request(t, "POST", "/user/login", body)
}

func userInfo(t *testing.T) {
	request(t, "GET", "/user/info", nil)
}
