// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/zhenorzz/goploy/controller"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/ws"
)

// Init router
func Init() {
	var rt = core.NewRouter()
	// rt.Middleware(example)
	rt.Add(ws.GetHub())
	rt.Add(controller.User{})
	rt.Add(controller.Namespace{})
	rt.Add(controller.Role{})
	rt.Add(controller.Project{})
	rt.Add(controller.Repository{})
	rt.Add(controller.Monitor{})
	rt.Add(controller.Deploy{})
	rt.Add(controller.Server{})
	rt.Add(controller.Log{})
	rt.Add(controller.Cron{})
	rt.Add(controller.Agent{})

	rt.Start()
}
