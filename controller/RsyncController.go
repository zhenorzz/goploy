package controller

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Rsync string

func (rsync *Rsync) Add(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("rsync", "-av", "192.168.21.146:/tmp", "/home/zane")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(out.String()))
}
