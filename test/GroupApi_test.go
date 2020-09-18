package test

import (
	"net/http"
	"testing"
)

func getGroupList(t *testing.T) {
	request(t, http.MethodGet, "/group/getList?page=1&rows=10", nil)
}

func getGroupOption(t *testing.T) {
	request(t, http.MethodGet, "/group/getOption", nil)
}

func getDeployOption(t *testing.T) {
	request(t, http.MethodGet, "/group/getDeployOption", nil)
}

func addGroup(t *testing.T) {
	body := struct {
		Name string `json:"name"`
	}{
		Name: getRandomStringOf(5),
	}
	resp := request(t, http.MethodPost, "/group/add", body)
	groupID = int64(resp.Data.(map[string]interface{})["id"].(float64))
}

func editGroup(t *testing.T) {
	body := struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}{
		ID:   groupID,
		Name: getRandomStringOf(5),
	}
	request(t, http.MethodPost, "/group/edit", body)
}

func removeGroup(t *testing.T) {
	body := struct {
		ID int64 `json:"id"`
	}{
		ID: groupID,
	}
	request(t, http.MethodDelete, "/group/remove", body)
}
