// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/vearutop/statigz"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/web"
	"io/fs"
	"log"
	"net/http"
	"net/url"
)

type Goploy struct {
	UserInfo  model.User
	Namespace struct {
		ID            int64
		PermissionIDs map[int64]struct{}
	}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	URLQuery       url.Values
	Body           []byte
}

type Response interface {
	Write(http.ResponseWriter, *http.Request) error
}

type Server struct {
	http.Server
	*Router
}

func (srv *Server) Spin() {
	if config.Toml.Env == "production" {
		subFS, err := fs.Sub(web.Dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/assets/", statigz.FileServer(subFS.(fs.ReadDirFS)))
		http.Handle("/favicon.ico", statigz.FileServer(subFS.(fs.ReadDirFS)))
	}
	http.Handle("/", srv.Router)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
