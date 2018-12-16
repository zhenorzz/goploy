package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id      string
	Balance uint64
}

func (user *User) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i'm in user controller")
	u := User{Id: "US123", Balance: 8}
	json.NewEncoder(w).Encode(u)
}
