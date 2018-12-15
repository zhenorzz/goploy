package controller

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type User struct{
	Id      string
	Balance uint64
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i'm in user controller")
	u := User{Id: "US123", Balance: 8}
	json.NewEncoder(w).Encode(u)
}