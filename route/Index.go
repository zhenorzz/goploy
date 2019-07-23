package route

import (
	"github.com/zhenorzz/goploy/controller"
	router "github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
)

// Init router
func Init() {
	var rt = new(router.Router)
	// rt.Middleware(exaplme)

	// home route
	rt.Add("/index/get", controller.Index{}.Get)
	// user route
	rt.Add("/user/login", controller.User{}.Login)
	rt.Add("/user/info", controller.User{}.Info)
	rt.Add("/user/getList", controller.User{}.GetList)
	rt.Add("/user/getOption", controller.User{}.GetOption)
	rt.Add("/user/add", controller.User{}.Add)
	rt.Add("/user/edit", controller.User{}.Edit)
	rt.Add("/user/remove", controller.User{}.Remove)
	rt.Add("/user/changePassword", controller.User{}.ChangePassword)

	// project route
	rt.Add("/project/getList", controller.Project{}.GetList)
	rt.Add("/project/getBindServerList", controller.Project{}.GetBindServerList)
	rt.Add("/project/getBindUserList", controller.Project{}.GetBindUserList)
	rt.Add("/project/add", controller.Project{}.Add)
	rt.Add("/project/edit", controller.Project{}.Edit)
	rt.Add("/project/remove", controller.Project{}.Remove)
	rt.Add("/project/addServer", controller.Project{}.AddServer)
	rt.Add("/project/addUser", controller.Project{}.AddUser)
	rt.Add("/project/removeProjectServer", controller.Project{}.RemoveProjectServer)
	rt.Add("/project/removeProjectUser", controller.Project{}.RemoveProjectUser)

	// deploy route
	rt.Add("/deploy/getList", controller.Deploy{}.GetList)
	rt.Add("/deploy/getDetail", controller.Deploy{}.GetDetail)
	rt.Add("/deploy/getSyncDetail", controller.Deploy{}.GetSyncDetail)
	rt.Add("/deploy/sync", controller.Deploy{}.Sync)
	rt.Add("/deploy/publish", controller.Deploy{}.Publish, middleware.HasPublishAuth)

	// server route
	rt.Add("/server/getList", controller.Server{}.GetList)
	rt.Add("/server/getOption", controller.Server{}.GetOption)
	rt.Add("/server/add", controller.Server{}.Add)
	rt.Add("/server/edit", controller.Server{}.Edit)
	rt.Add("/server/remove", controller.Server{}.Remove)

	// role route
	rt.Add("/role/getOption", controller.Role{}.GetOption)
	rt.Add("/role/getPermissionList", controller.Role{}.GetPermissionList)
	rt.Add("/role/edit", controller.Role{}.Edit)
	rt.Start()
}
