package test

import (
	"net/http"
	"testing"
)

func getTemplateList(t *testing.T) {
	request(t, http.MethodGet, "/template/getList?page=1&rows=10", nil)
}

func getTemplateOption(t *testing.T) {
	request(t, http.MethodGet, "/template/getOption", nil)
}

func addTemplate(t *testing.T) {
	body := struct {
		Name         string `json:"name"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script"`
	}{
		Name:   getRandomStringOf(5),
		Script: "echo 1",
	}
	resp := request(t, http.MethodPost, "/template/add", body)
	templateID = int64(resp.Data.(map[string]interface{})["id"].(float64))
}

func editTemplate(t *testing.T) {
	body := struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Remark       string `json:"remark"`
		PackageIDStr string `json:"packageIdStr"`
		Script       string `json:"script"`
	}{
		ID:     templateID,
		Name:   getRandomStringOf(5),
		Script: "echo 2",
	}
	request(t, http.MethodPost, "/template/edit", body)
}

func removeTemplate(t *testing.T) {
	body := struct {
		ID int64 `json:"id"`
	}{
		ID: templateID,
	}
	request(t, http.MethodDelete, "/template/remove", body)
}
