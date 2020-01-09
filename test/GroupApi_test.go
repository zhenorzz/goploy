package test

import (
	router "goploy/core"
	"testing"
)

func getGroupList(t *testing.T) {
	request(t, router.GET, "/group/getList?page=1&rows=10", nil)
}

func getGroupOption(t *testing.T) {
	request(t, router.GET, "/group/getOption", nil)
}

func getDeployOption(t *testing.T) {
	request(t, router.GET, "/group/getDeployOption", nil)
}

func addGroup(t *testing.T) {
	body := struct {
		Name string `json:"name"`
	}{
		Name: getRandomStringOf(5),
	}
	resp := request(t, router.POST, "/group/add", body)
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
	request(t, router.POST, "/group/edit", body)
}

func removeGroup(t *testing.T) {
	body := struct {
		ID int64 `json:"id"`
	}{
		ID: groupID,
	}
	request(t, router.DELETE, "/group/remove", body)
}
