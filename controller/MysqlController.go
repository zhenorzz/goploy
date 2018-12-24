package controller

import (
	"fmt"
	"net/http"

	"github.com/zhenorzz/goploy/model"
)

type Mysql string

func (mysql *Mysql) Query(w http.ResponseWriter, r *http.Request) {
	ceshis := new(model.Ceshis)
	ceshis.QueryMany()
	fmt.Println(ceshis)
}
