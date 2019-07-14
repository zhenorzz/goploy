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
	rt.Add("/user/getList", new(controller.User).GetList)
	rt.Add("/user/getOption", new(controller.User).GetOption)
	rt.Add("/user/add", new(controller.User).Add)
	rt.Add("/user/edit", new(controller.User).Edit)
	rt.Add("/user/changePassword", new(controller.User).ChangePassword)

	// project route
	rt.Add("/project/getList", new(controller.Project).GetList)
	rt.Add("/project/getBindServerList", new(controller.Project).GetBindServerList)
	rt.Add("/project/getBindUserList", new(controller.Project).GetBindUserList)
	rt.Add("/project/create", new(controller.Project).Create)
	rt.Add("/project/add", new(controller.Project).Add)
	rt.Add("/project/edit", new(controller.Project).Edit)
	rt.Add("/project/addServer", new(controller.Project).AddServer)
	rt.Add("/project/addUser", new(controller.Project).AddUser)
	rt.Add("/project/removeProjectServer", new(controller.Project).RemoveProjectServer)
	rt.Add("/project/removeProjectUser", new(controller.Project).RemoveProjectUser)

	// deploy route
	rt.Add("/deploy/getList", new(controller.Deploy).GetList)
	rt.Add("/deploy/getDetail", new(controller.Deploy).GetDetail)
	rt.Add("/deploy/publish", new(controller.Deploy).Publish)

	// server route
	rt.Add("/server/getList", new(controller.Server).GetList)
	rt.Add("/server/getOption", new(controller.Server).GetOption)
	rt.Add("/server/add", new(controller.Server).Add)
	rt.Add("/server/edit", new(controller.Server).Edit)

	// role route
	rt.Add("/role/getOption", new(controller.Role).GetOption)
	rt.Start()
}
