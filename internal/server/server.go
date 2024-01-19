// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import (
	"errors"
	"github.com/vearutop/statigz"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/web"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response interface {
	Write(http.ResponseWriter, *http.Request) error
}

type Server struct {
	http.Server
	*Router
}

func (srv *Server) Spin() {
	pid := strconv.Itoa(os.Getpid())
	_ = os.WriteFile(config.GetPidFile(), []byte(pid), 0755)
	println("Started at " + time.Now().String())
	println("goploy -h for more help")
	println("Current pid:   " + pid)
	println("Config Loaded: " + config.GetConfigFile())
	println("Env:           " + config.Toml.Env)
	println("Log:           " + config.Toml.Log.Path)
	println("Listen:        " + config.Toml.Web.Port)

	if config.Toml.Env == "production" {
		subFS, err := fs.Sub(web.Dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/assets/", statigz.FileServer(subFS.(fs.ReadDirFS)))
		http.Handle("/favicon.ico", statigz.FileServer(subFS.(fs.ReadDirFS)))
	}
	http.Handle("/", srv.Router)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("ListenAndServe: ", err.Error())
	}

	_ = os.Remove(config.GetPidFile())
	println("Shutdown success")
}
