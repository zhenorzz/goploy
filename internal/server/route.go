// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import "errors"

type Route struct {
	pattern       string
	method        string                    // Method specifies the HTTP method (GET, POST, PUT, etc.).
	permissionIDs []int64                   // permission list
	white         bool                      // no need to login
	middlewares   []func(gp *Goploy) error  // Middlewares run before callback, trigger error will end the request
	callback      func(gp *Goploy) Response // API function
	logFunc       func(gp *Goploy, resp Response)
}

type RouteHandler interface {
	Handler() []Route
}

func NewRoute(pattern, method string, callback func(gp *Goploy) Response) Route {
	return newRoute(pattern, method, callback)
}

func NewWhiteRoute(pattern, method string, callback func(gp *Goploy) Response) Route {
	route := newRoute(pattern, method, callback)
	route.white = true
	return route
}

func newRoute(pattern, method string, callback func(gp *Goploy) Response) Route {
	return Route{
		pattern:  pattern,
		method:   method,
		callback: callback,
	}
}

func (r Route) Permissions(permissionIDs ...int64) Route {
	for _, permissionID := range permissionIDs {
		r.permissionIDs = append(r.permissionIDs, permissionID)
	}
	return r
}

// Middleware global Middleware handle function
func (r Route) Middleware(middleware func(gp *Goploy) error) Route {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// LogFunc callback finished
func (r Route) LogFunc(f func(gp *Goploy, resp Response)) Route {
	r.logFunc = f
	return r
}

func (r Route) hasPermission(permissionIDs map[int64]struct{}) error {
	if len(r.permissionIDs) == 0 {
		return nil
	}

	for _, permissionID := range r.permissionIDs {
		if _, ok := permissionIDs[permissionID]; ok {
			return nil
		}
	}

	return errors.New("no permission")
}
