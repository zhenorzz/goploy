package route

import (
	"github.com/zhenorzz/goploy/controller"
	router "github.com/zhenorzz/goploy/core"
)

// Init router
func Init() {
	var rt = new(router.Router)
	rt.Middleware(CheckToken)
	rt.Add("/user/login", new(controller.User).Login)
	rt.Add("/user/info", new(controller.User).Info)

	// project route
	rt.Add("/project/add", new(controller.Project).Add)
	rt.Add("/project/get", new(controller.Project).Get)
	rt.Add("/project/create", new(controller.Project).Create)
	rt.Add("/project/branch", new(controller.Project).Branch)
	rt.Add("/project/commit", new(controller.Project).Commit)

	// deploy route
	rt.Add("/deploy/add", new(controller.Deploy).Add)
	rt.Add("/deploy/get", new(controller.Deploy).Get)

	rt.Add("/github/search", new(controller.Github).Search)
	rt.Add("/rsync/add", new(controller.Rsync).Add)
	rt.Add("/mysql/query", new(controller.Mysql).Query)
	rt.Start()
}
