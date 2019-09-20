package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"goploy/core"
	"goploy/model"
	"goploy/route"
	"goploy/ws"
)

func main() {
	godotenv.Load(core.GolbalPath + ".env")
	core.CreateValidator()
	model.Init()
	ws.Init()
	route.Init()
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
