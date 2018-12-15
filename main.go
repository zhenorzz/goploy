package main

import (
	"net/http"
	"log"
	"github.com/zhenorzz/goploy/route"
	_ "github.com/joho/godotenv"
)

func main() {
	route.Init()
	err := http.ListenAndServe(":9091", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}