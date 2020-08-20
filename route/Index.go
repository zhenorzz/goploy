package route

import (
	"github.com/zhenorzz/goploy/controller"
	"github.com/zhenorzz/goploy/core"
	router "github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/ws"
)

// Init router
func Init() *router.Router {
	var rt = new(router.Router)
	// rt.Middleware(example)
	// no need to check login
	rt.RegisterWhiteList(map[string]struct{}{
		"/user/login":        {},
		"/user/isShowPhrase": {},
		"/deploy/webhook":    {},
	})
	// websocket route
	rt.Add("/ws/connect", router.GET, ws.GetHub().Connect)

	// user route
	rt.Add("/user/login", router.POST, controller.User{}.Login)
	rt.Add("/user/info", router.GET, controller.User{}.Info)
	rt.Add("/user/getList", router.GET, controller.User{}.GetList)
	rt.Add("/user/getTotal", router.GET, controller.User{}.GetTotal)
	rt.Add("/user/getOption", router.GET, controller.User{}.GetOption)
	rt.Add("/user/add", router.POST, controller.User{}.Add).Role(core.RoleAdmin)
	rt.Add("/user/edit", router.POST, controller.User{}.Edit).Role(core.RoleAdmin)
	rt.Add("/user/remove", router.DELETE, controller.User{}.Remove).Role(core.RoleAdmin)
	rt.Add("/user/changePassword", router.POST, controller.User{}.ChangePassword)

	// namespace route
	rt.Add("/namespace/getList", router.GET, controller.Namespace{}.GetList)
	rt.Add("/namespace/getTotal", router.GET, controller.Namespace{}.GetTotal)
	rt.Add("/namespace/getBindUserList", router.GET, controller.Namespace{}.GetBindUserList)
	rt.Add("/namespace/getUserOption", router.GET, controller.Namespace{}.GetUserOption)
	rt.Add("/namespace/add", router.POST, controller.Namespace{}.Add).Role(core.RoleAdmin)
	rt.Add("/namespace/edit", router.POST, controller.Namespace{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/addUser", router.POST, controller.Namespace{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/removeUser", router.DELETE, controller.Namespace{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager})

	// project route
	rt.Add("/project/getList", router.GET, controller.Project{}.GetList)
	rt.Add("/project/getTotal", router.GET, controller.Project{}.GetTotal)
	rt.Add("/project/getRemoteBranchList", router.GET, controller.Project{}.GetRemoteBranchList)
	rt.Add("/project/getBindServerList", router.GET, controller.Project{}.GetBindServerList)
	rt.Add("/project/getBindUserList", router.GET, controller.Project{}.GetBindUserList)
	rt.Add("/project/add", router.POST, controller.Project{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/edit", router.POST, controller.Project{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/setAutoDeploy", router.POST, controller.Project{}.SetAutoDeploy).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/remove", router.DELETE, controller.Project{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addServer", router.POST, controller.Project{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addUser", router.POST, controller.Project{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeServer", router.DELETE, controller.Project{}.RemoveServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeUser", router.DELETE, controller.Project{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addTask", router.POST, controller.Project{}.AddTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/editTask", router.POST, controller.Project{}.EditTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeTask", router.POST, controller.Project{}.RemoveTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/getTaskList", router.GET, controller.Project{}.GetTaskList).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	// monitor route
	rt.Add("/monitor/getList", router.GET, controller.Monitor{}.GetList)
	rt.Add("/monitor/getTotal", router.GET, controller.Monitor{}.GetTotal)
	rt.Add("/monitor/check", router.POST, controller.Monitor{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/add", router.POST, controller.Monitor{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/edit", router.POST, controller.Monitor{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/toggle", router.POST, controller.Monitor{}.Toggle).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/remove", router.DELETE, controller.Monitor{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	//// deploy route
	rt.Add("/deploy/getList", router.GET, controller.Deploy{}.GetList)
	rt.Add("/deploy/getDetail", router.GET, controller.Deploy{}.GetDetail)
	rt.Add("/deploy/getCommitList", router.GET, controller.Deploy{}.GetCommitList)
	rt.Add("/deploy/getPreview", router.GET, controller.Deploy{}.GetPreview)
	rt.Add("/deploy/publish", router.POST, controller.Deploy{}.Publish, middleware.HasPublishAuth)
	rt.Add("/deploy/webhook", router.POST, controller.Deploy{}.Webhook, middleware.FilterEvent)

	// server route
	rt.Add("/server/getList", router.GET, controller.Server{}.GetList)
	rt.Add("/server/getTotal", router.GET, controller.Server{}.GetTotal)
	rt.Add("/server/getInstallPreview", router.GET, controller.Server{}.GetInstallPreview)
	rt.Add("/server/getInstallList", router.GET, controller.Server{}.GetInstallList)
	rt.Add("/server/getOption", router.GET, controller.Server{}.GetOption)
	rt.Add("/server/check", router.POST, controller.Server{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/add", router.POST, controller.Server{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/edit", router.POST, controller.Server{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/remove", router.DELETE, controller.Server{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/install", router.POST, controller.Server{}.Install).Roles([]string{core.RoleAdmin, core.RoleManager})

	// template route
	rt.Add("/template/getList", router.GET, controller.Template{}.GetList)
	rt.Add("/template/getTotal", router.GET, controller.Template{}.GetTotal)
	rt.Add("/template/getOption", router.GET, controller.Template{}.GetOption)
	rt.Add("/template/add", router.POST, controller.Template{}.Add)
	rt.Add("/template/edit", router.POST, controller.Template{}.Edit)
	rt.Add("/template/remove", router.DELETE, controller.Template{}.Remove)

	// template route
	rt.Add("/package/getList", router.GET, controller.Package{}.GetList)
	rt.Add("/package/getTotal", router.GET, controller.Package{}.GetTotal)
	rt.Add("/package/getOption", router.GET, controller.Package{}.GetOption)
	rt.Add("/package/upload", router.POST, controller.Package{}.Upload)

	// crontab route
	rt.Add("/crontab/getList", router.GET, controller.Crontab{}.GetList)
	rt.Add("/crontab/getTotal", router.GET, controller.Crontab{}.GetTotal)
	rt.Add("/crontab/getRemoteServerList", router.GET, controller.Crontab{}.GetRemoteServerList)
	rt.Add("/crontab/getBindServerList", router.GET, controller.Crontab{}.GetBindServerList)
	rt.Add("/crontab/add", router.POST, controller.Crontab{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/edit", router.POST, controller.Crontab{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/import", router.POST, controller.Crontab{}.Import).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/remove", router.DELETE, controller.Crontab{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/addServer", router.POST, controller.Crontab{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/removeCrontabServer", router.DELETE, controller.Crontab{}.RemoveCrontabServer).Roles([]string{core.RoleAdmin, core.RoleManager})

	rt.Start()
	return rt
}
