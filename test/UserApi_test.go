package test

//
//import (
//	"net/http"
//	"testing"
//)
//
//func userLogin(t *testing.T) {
//	body := struct {
//		Account  string `json:"account"`
//		Password string `json:"password"`
//	}{
//		Account:  "admin",
//		Password: "admin!@#",
//	}
//	resp := request(t, http.MethodPost, "/user/login", body)
//	token = resp.Data.(map[string]interface{})["token"].(string)
//}
//
//func userInfo(t *testing.T) {
//	request(t, http.MethodGet, "/user/info", nil)
//}
//
//func getUserList(t *testing.T) {
//	request(t, http.MethodGet, "/user/getList?page=1&rows=10", nil)
//}
//
//func getUserOption(t *testing.T) {
//	request(t, http.MethodGet, "/user/getOption", nil)
//}
//
//func addUser(t *testing.T) {
//	body := struct {
//		Account        string `json:"account"`
//		Password       string `json:"password"`
//		Name           string `json:"name"`
//		Mobile         string `json:"mobile"`
//		Role           string `json:"role"`
//		ManageGroupStr string `json:"manageGroupStr"`
//	}{
//		Account:        getRandomStringOf(5),
//		Password:       "admin!@#",
//		Name:           "name",
//		Role:           "admin",
//		ManageGroupStr: "all",
//	}
//	resp := request(t, http.MethodPost, "/user/add", body)
//	userID = int64(resp.Data.(map[string]interface{})["id"].(float64))
//}
//
//func editUser(t *testing.T) {
//	body := struct {
//		ID             int64  `json:"id"`
//		Password       string `json:"password"`
//		Name           string `json:"name"`
//		Mobile         string `json:"mobile"`
//		Role           string `json:"role"`
//		ManageGroupStr string `json:"manageGroupStr"`
//	}{
//		ID:             userID,
//		Password:       "admin!@#",
//		Name:           "change_name",
//		Mobile:         "13800138000",
//		Role:           "admin",
//		ManageGroupStr: "all",
//	}
//	request(t, http.MethodPost, "/user/edit", body)
//}
//
//func removeUser(t *testing.T) {
//	body := struct {
//		ID int64 `json:"id"`
//	}{
//		ID: userID,
//	}
//	request(t, http.MethodDelete, "/user/remove", body)
//}
//
//func changeUserPassword(t *testing.T) {
//	body := struct {
//		OldPassword string `json:"oldPwd"`
//		NewPassword string `json:"newPwd"`
//	}{
//		OldPassword: "admin!@#",
//		NewPassword: "admin!@#",
//	}
//	request(t, http.MethodPost, "/user/changePassword", body)
//}
