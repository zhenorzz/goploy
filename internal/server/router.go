// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/web"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Router is Route slice and global middlewares
type Router struct {
	routes      map[string]Route
	middlewares *[]func(gp *Goploy) error // Middlewares run before all Route
}

func NewRouter() *Router {
	router := Router{
		routes:      map[string]Route{},
		middlewares: new([]func(gp *Goploy) error),
	}

	return &router
}

func (rt *Router) Middleware(middleware func(gp *Goploy) error) {
	*rt.middlewares = append(*rt.middlewares, middleware)
}

func (rt *Router) Register(ra RouteHandler) {
	for _, r := range ra.Handler() {
		rt.routes[r.pattern] = r
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if config.Toml.CORS.Enabled {
		w.Header().Set("Access-Control-Allow-Origin", config.Toml.CORS.Origins)
		w.Header().Add("Access-Control-Allow-Headers", config.Toml.CORS.Headers)
		w.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(config.Toml.CORS.Credentials))
		w.Header().Add("Access-Control-Allow-Methods", config.Toml.CORS.Methods)
	}

	// If in production env, serve file in go server,
	// else serve file in npm
	if config.Toml.Env == "production" {
		if "/" == r.URL.Path {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			index, err := web.Dist.Open("dist/index.html")
			if err != nil {
				fmt.Fprint(w, "404")
				log.Error(err.Error())
			}
			defer index.Close()
			contents, err := io.ReadAll(index)
			fmt.Fprint(w, string(contents))
			return
		}
	}

	_, resp := rt.doRequest(w, r)
	if err := resp.Write(w, r); err != nil {
		log.Error(err.Error())
	}
	return
}

func (rt *Router) doRequest(w http.ResponseWriter, r *http.Request) (*Goploy, Response) {
	gp := new(Goploy)
	route, ok := rt.routes[r.URL.Path]
	if !ok {
		return gp, response.JSON{Code: response.Deny, Message: "No such method"}
	}
	if route.method != r.Method {
		return gp, response.JSON{Code: response.IllegalRequest, Message: "Invalid request method"}
	}

	if !route.white {
		namespaceIDRaw := r.Header.Get(config.NamespaceHeaderName)
		if namespaceIDRaw == "" {
			namespaceIDRaw = r.URL.Query().Get(config.NamespaceHeaderName)
		}
		var namespaceID int64
		if namespaceIDRaw != "" {
			_id, err := strconv.ParseInt(namespaceIDRaw, 10, 64)
			if err != nil {
				return gp, response.JSON{Code: response.Deny, Message: "Invalid namespace"}
			}
			namespaceID = _id
		}

		apiKey := r.Header.Get(config.ApiKeyHeaderName)
		if apiKey == "" {
			apiKey = r.URL.Query().Get(config.ApiKeyHeaderName)
		}

		unParseToken := ""
		// check token
		goployTokenCookie, err := r.Cookie(config.Toml.Cookie.Name)
		if err != nil {
			unParseToken = r.URL.Query().Get(config.Toml.Cookie.Name)
		} else {
			unParseToken = goployTokenCookie.Value
		}

		if unParseToken != "" {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(unParseToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Toml.JWT.Key), nil
			})

			if err != nil || !token.Valid {
				return gp, response.JSON{Code: response.LoginExpired, Message: "Login expired"}
			}

			gp.UserInfo, err = model.User{ID: int64(claims["id"].(float64))}.GetData()
			if err != nil {
				return gp, response.JSON{Code: response.Deny, Message: "Get user information error"}
			}
		} else if apiKey != "" {
			gp.UserInfo, err = model.User{ApiKey: apiKey}.GetDataByApiKey()
			if err != nil {
				return gp, response.JSON{Code: response.Deny, Message: "Get user information using api key error"}
			}
		} else {
			return gp, response.JSON{Code: response.IllegalRequest, Message: "Illegal request"}
		}

		if gp.UserInfo.State != 1 {
			return gp, response.JSON{Code: response.AccountDisabled, Message: "No available user"}
		}

		if gp.UserInfo.SuperManager == model.SuperManager {
			permissionIDs, err := model.Permission{}.GetIDs()
			if err != nil {
				return gp, response.JSON{Code: response.Deny, Message: err.Error()}
			}
			gp.Namespace.ID = namespaceID
			gp.Namespace.PermissionIDs = permissionIDs
		} else if namespaceID > 0 {
			namespace, err := model.NamespaceUser{
				NamespaceID: namespaceID,
				UserID:      gp.UserInfo.ID,
			}.GetDataByUserNamespace()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return gp, response.JSON{Code: response.NamespaceInvalid, Message: "No available namespace"}
				} else {
					return gp, response.JSON{Code: response.Deny, Message: err.Error()}
				}
			}
			gp.Namespace.ID = namespaceID
			gp.Namespace.PermissionIDs = namespace.PermissionIDs
		} else {
			gp.Namespace.PermissionIDs = map[int64]struct{}{}
		}

		if err = route.hasPermission(gp.Namespace.PermissionIDs); err != nil {
			return gp, response.JSON{Code: response.Deny, Message: err.Error()}
		}

		goployTokenStr, err := gp.UserInfo.CreateToken()
		if err == nil {
			// update jwt time
			cookie := http.Cookie{Name: config.Toml.Cookie.Name, Value: goployTokenStr, Path: "/", MaxAge: config.Toml.Cookie.Expire, HttpOnly: true}
			http.SetCookie(w, &cookie)
		}
	}
	r.Header.Add("_time", time.Now().Format("20060102150405"))
	gp.Request = r
	gp.ResponseWriter = w
	gp.URLQuery = r.URL.Query()

	// save the body request data because io.ReadAll will clear the requestBody
	if r.ContentLength > 0 && hasContentType(r, "application/json") {
		gp.Body, _ = io.ReadAll(r.Body)
	}

	// common middlewares
	for _, middleware := range *rt.middlewares {
		err := middleware(gp)
		if err != nil {
			return gp, response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	// route middlewares
	for _, middleware := range route.middlewares {
		if err := middleware(gp); err != nil {
			return gp, response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	resp := route.callback(gp)

	if route.logFunc != nil {
		route.logFunc(gp, resp)
	}

	return gp, resp
}

func hasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return false
	}
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
