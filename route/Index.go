package route

import (
	"github.com/zhenorzz/goploy/controller"
	"github.com/zhenorzz/goploy/core"
	router "github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/ws"
)

// Init router
func Init() {
	var rt = new(router.Router)
	// rt.Middleware(exaplme)

	// home route
	rt.Add("/index/get", controller.Index{}.Get)

	// websocket route
	rt.Add("/ws/unicast", ws.GetUnicastHub().Unicast)

	// user route
	rt.Add("/user/login", controller.User{}.Login)
	rt.Add("/user/info", controller.User{}.Info)
	rt.Add("/user/getList", controller.User{}.GetList)
	rt.Add("/user/getOption", controller.User{}.GetOption)
	rt.Add("/user/add", controller.User{}.Add).Role(core.RoleAdmin)
	rt.Add("/user/edit", controller.User{}.Edit).Role(core.RoleAdmin)
	rt.Add("/user/remove", controller.User{}.Remove).Role(core.RoleAdmin)
	rt.Add("/user/changePassword", controller.User{}.ChangePassword)

	// project route
	rt.Add("/project/getList", controller.Project{}.GetList)
	rt.Add("/project/getBindServerList", controller.Project{}.GetBindServerList)
	rt.Add("/project/getBindUserList", controller.Project{}.GetBindUserList)
	rt.Add("/project/add", controller.Project{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/edit", controller.Project{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/remove", controller.Project{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addServer", controller.Project{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addUser", controller.Project{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeProjectServer", controller.Project{}.RemoveProjectServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeProjectUser", controller.Project{}.RemoveProjectUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	// deploy route
	rt.Add("/deploy/getList", controller.Deploy{}.GetList)
	rt.Add("/deploy/getDetail", controller.Deploy{}.GetDetail)
	rt.Add("/deploy/getCommitList", controller.Deploy{}.GetCommitList)
	rt.Add("/deploy/getPreview", controller.Deploy{}.GetPreview)
	rt.Add("/deploy/publish", controller.Deploy{}.Publish, middleware.HasPublishAuth)
	rt.Add("/deploy/rollback", controller.Deploy{}.Rollback, middleware.HasPublishAuth)

	// server route
	rt.Add("/server/getList", controller.Server{}.GetList)
	rt.Add("/server/getInstallPreview", controller.Server{}.GetInstallPreview)
	rt.Add("/server/getInstallList", controller.Server{}.GetInstallList)
	rt.Add("/server/getOption", controller.Server{}.GetOption)
	rt.Add("/server/add", controller.Server{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/server/edit", controller.Server{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/server/remove", controller.Server{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/server/install", controller.Server{}.Install).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	// template route
	rt.Add("/template/getList", controller.Template{}.GetList)
	rt.Add("/template/getOption", controller.Template{}.GetOption)
	rt.Add("/template/add", controller.Template{}.Add)
	rt.Add("/template/edit", controller.Template{}.Edit)
	rt.Add("/template/remove", controller.Template{}.Remove)

	// template route
	rt.Add("/package/getList", controller.Package{}.GetList)
	rt.Add("/package/getOption", controller.Package{}.GetOption)
	rt.Add("/package/upload", controller.Package{}.Upload)

	// projectGroup route
	rt.Add("/group/getList", controller.Group{}.GetList)
	rt.Add("/group/getOption", controller.Group{}.GetOption)
	rt.Add("/group/add", controller.Group{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/group/edit", controller.Group{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/group/remove", controller.Group{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})

	// role route
	rt.Add("/role/getOption", controller.Role{}.GetOption)
	rt.Start()
}
