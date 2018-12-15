package controller

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type Commits struct {
	Sha    string `json:"sha"`
	NodeId string `json:"node_id"`
	Commit Commit `json:"commit"`
}
type Commit struct {
	Committer Committer `json:"committer"`
	Message   string    `json:"message"`
	Tree      Tree      `json:"tree"`
}
type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}
type Tree struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

func GithubSearch(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/godis/commits/btree")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var commit Commits
	json.Unmarshal(body, &commit)
	fmt.Println(commit)
	w.Write(body)
}
