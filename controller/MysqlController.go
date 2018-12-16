package controller

import (
	"net/http"

	"github.com/zhenorzz/goploy/model"
)

func MysqlQuery(w http.ResponseWriter, r *http.Request) {
	model.Query()
}
