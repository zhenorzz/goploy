package test

import (
	"net/http"
	"strconv"
	"testing"
)

func getServerList(t *testing.T) {
	request(t, http.MethodGet, "/server/getList?page=1&rows=10", nil)
}

func GetServerInstallPreview(t *testing.T) {
	request(t, http.MethodGet, "/server/getInstallPreview?serverId="+strconv.FormatInt(serverID, 10), nil)
}

func getServerInstallList(t *testing.T) {
	request(t, http.MethodGet, "/server/getInstallList?token=5861dc82-061d-401d-9152-ba7a2527edf7", nil)
}

func getServerOption(t *testing.T) {
	request(t, http.MethodGet, "/server/getOption", nil)
}

func addServer(t *testing.T) {
	body := struct {
		Name    string `json:"name"`
		IP      string `json:"ip"`
		Port    int    `json:"port"`
		Owner   string `json:"owner"`
		GroupID int64  `json:"groupId"`
	}{
		Name:    getRandomStringOf(5),
		IP:      "129.204.80.253",
		Port:    22,
		Owner:   "root",
		GroupID: 0,
	}
	resp := request(t, http.MethodPost, "/server/add", body)
	serverID = int64(resp.Data.(map[string]interface{})["id"].(float64))
}

func editServer(t *testing.T) {
	body := struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		IP      string `json:"ip"`
		Port    int    `json:"port"`
		Owner   string `json:"owner"`
		GroupID int64  `json:"groupId"`
	}{
		ID:      serverID,
		Name:    getRandomStringOf(5),
		IP:      "129.204.80.253",
		Port:    22,
		Owner:   "root",
		GroupID: 0,
	}
	request(t, http.MethodPost, "/server/edit", body)
}

func removeServer(t *testing.T) {
	body := struct {
		ID int64 `json:"id"`
	}{
		ID: serverID,
	}
	request(t, http.MethodDelete, "/server/remove", body)
}
