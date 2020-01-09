package test

import (
	router "goploy/core"
	"testing"
)

func getPackageList(t *testing.T) {
	request(t, router.GET, "/package/getList?page=1&rows=10", nil)
}

func getPackageOption(t *testing.T) {
	request(t, router.GET, "/package/getOption", nil)
}
