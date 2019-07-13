package route

import (
	"github.com/zhenorzz/goploy/controller"
	router "github.com/zhenorzz/goploy/core"
)

// Init router
func Init() {
	var rt = new(router.Router)
	rt.Middleware(CheckToken)

	// home route
	rt.Add("/index/get", new(controller.Index).Get)

	// user route
	rt.Add("/user/isShowPhrase", new(controller.User).IsShowPhrase)
	rt.Add("/user/login", new(controller.User).Login)
	rt.Add("/user/info", new(controller.User).Info)
	rt.Add("/user/get", new(controller.User).Get)
	rt.Add("/user/add", new(controller.User).Add)
	rt.Add("/user/changePassword", new(controller.User).ChangePassword)

	// project route
	rt.Add("/project/get", new(controller.Project).Get)
	rt.Add("/project/getDetail", new(controller.Project).GetDetail)
	rt.Add("/project/create", new(controller.Project).Create)
	rt.Add("/project/add", new(controller.Project).Add)

	// deploy route
	rt.Add("/deploy/get", new(controller.Deploy).Get)
	rt.Add("/deploy/publish", new(controller.Deploy).Publish)

	// server route
	rt.Add("/server/get", new(controller.Server).Get)
	rt.Add("/server/add", new(controller.Server).Add)
	rt.Add("/server/edit", new(controller.Server).Edit)
	rt.Start()
}
