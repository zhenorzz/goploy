package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv"
	"github.com/zhenorzz/goploy/route"
)

func main() {
	route.Init()
	err := http.ListenAndServe(":9091", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
