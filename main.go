package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/cron"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/route"
	"github.com/zhenorzz/goploy/ws"
)

func main() {
	godotenv.Load(core.GolbalPath + ".env")
	model.Init()
	go cron.Run()
	syncHub := ws.GetSyncHub()
	go syncHub.Run()
	route.Init()
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
