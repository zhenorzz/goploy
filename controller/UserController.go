package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id uint32
}

func (user *User) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i'm in user controller")
	u := User{Id: 1}
	json.NewEncoder(w).Encode(u)
}

func (user *User) Info(w http.ResponseWriter, r *http.Request) {
	u := User{Id: 1}
	json.NewEncoder(w).Encode(u)
}
