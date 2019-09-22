package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goploy/core"
	"goploy/model"
	"goploy/route"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestApi(t *testing.T) {
	godotenv.Load("../.env")
	core.CreateValidator()
	model.Init()
	// user login have token
	userLogin(t)
	addUser(t)
	t.Run("user/info", userInfo)
	t.Run("user/getList", getUserList)
	t.Run("user/getOption", getUserOption)
	t.Run("user/edit", editUser)
	t.Run("user/remove", removeUser)
	t.Run("user/changePassword", changeUserPassword)
}

var handler = route.Init()
var (
	token string
	userID int64
)

func request(t *testing.T, method, url string, body interface{}) core.Response {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(body)
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
	return resp
}

func getRandomStringOf(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bs := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(len(bs))])
	}
	return string(result)
}
