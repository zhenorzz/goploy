package test

import (
	"net/http"
	"testing"
)

func getPackageList(t *testing.T) {
	request(t, http.MethodGet, "/package/getList?page=1&rows=10", nil)
}

func getPackageOption(t *testing.T) {
	request(t, http.MethodGet, "/package/getOption", nil)
}
