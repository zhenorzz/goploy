package test

import (
	router "goploy/core"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func userLogin(t *testing.T) {
	//创建一个请求
	body := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}{
		Account:  "admin",
		Password: "admin!@#",
	}
	request(t, router.POST, "/user/login", body)
}

func userInfo(t *testing.T) {
	request(t, router.GET, "/user/info", nil)
}

func getUserList(t *testing.T) {
	request(t, router.GET, "/user/getList", nil)
}

func getUserOption(t *testing.T) {
	request(t, router.GET, "/user/getOption", nil)
}

func addUser(t *testing.T) {
	request(t, router.POST, "/user/add", nil)
}

func editUser(t *testing.T) {
	request(t, router.POST, "/user/edit", nil)
}

func removeUser(t *testing.T) {
	request(t, router.DELETE, "/user/remove", nil)
}

func changeUserPassword(t *testing.T) {
	request(t, router.POST, "/user/changePassword", nil)
}
