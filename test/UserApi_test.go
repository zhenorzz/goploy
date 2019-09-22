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
	resp := request(t, router.POST, "/user/login", body)
	token = resp.Data.(map[string]interface{})["token"].(string)
}

func userInfo(t *testing.T) {
	request(t, router.GET, "/user/info", nil)
}

func getUserList(t *testing.T) {
	request(t, router.GET, "/user/getList?page=1&rows=10", nil)
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
		Account:        getRandomStringOf(5),
		Password:       "admin!@#",
		Name:           "name",
		Role:           "admin",
		ManageGroupStr: "all",
	}
	resp := request(t, router.POST, "/user/add", body)
	userID = int64(resp.Data.(map[string]interface{})["id"].(float64))
}

func editUser(t *testing.T) {
	body := struct {
		ID             int64  `json:"id"`
		Password       string `json:"password"`
		Name           string `json:"name"`
		Mobile         string `json:"mobile"`
		Role           string `json:"role"`
		ManageGroupStr string `json:"manageGroupStr"`
	}{
		ID:             userID,
		Password:       "admin!@#",
		Name:           "change_name",
		Mobile:         "13800138000",
		Role:           "admin",
		ManageGroupStr: "all",
	}
	request(t, router.POST, "/user/edit", body)
}

func removeUser(t *testing.T) {
	request(t, router.DELETE, "/user/remove", nil)
}

func changeUserPassword(t *testing.T) {
	request(t, router.POST, "/user/changePassword", nil)
}
