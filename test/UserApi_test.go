package test

import (
	router "goploy/core"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func userLogin(t *testing.T) {
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
	body := struct {
		Account        string `json:"account"`
		Password       string `json:"password"`
		Name           string `json:"name"`
		Mobile         string `json:"mobile"`
		Role           string `json:"role"`
		ManageGroupStr string `json:"manageGroupStr"`
	}{
		Account:  "admin",
		Password: "admin!@#",
		Name:     "name",
		Role:     "admin",
	}
	request(t, router.POST, "/user/add", body)
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
