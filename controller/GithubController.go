package controller

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

func GithubSearch(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/godis/commits/btree")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
