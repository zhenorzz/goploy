package route

import (
	"github.com/zhenorzz/goploy/controller"
	"github.com/zhenorzz/goploy/core"
	router "github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/ws"
	"net/http"
)

// Init router
func Init() *router.Router {
	var rt = new(router.Router)
	// rt.Middleware(example)
	// no need to check login
	rt.RegisterWhiteList(map[string]struct{}{
		"/user/login":      {},
		"/deploy/webhook":  {},
		"/deploy/callback": {},
	})
	// websocket route
	rt.Add("/ws/connect", http.MethodGet, ws.GetHub().Connect)

	// user route
	rt.Add("/user/login", http.MethodPost, controller.User{}.Login)
	rt.Add("/user/info", http.MethodGet, controller.User{}.Info)
	rt.Add("/user/getList", http.MethodGet, controller.User{}.GetList)
	rt.Add("/user/getTotal", http.MethodGet, controller.User{}.GetTotal)
	rt.Add("/user/getOption", http.MethodGet, controller.User{}.GetOption)
	rt.Add("/user/add", http.MethodPost, controller.User{}.Add).Role(core.RoleAdmin)
	rt.Add("/user/edit", http.MethodPost, controller.User{}.Edit).Role(core.RoleAdmin)
	rt.Add("/user/remove", http.MethodDelete, controller.User{}.Remove).Role(core.RoleAdmin)
	rt.Add("/user/changePassword", http.MethodPost, controller.User{}.ChangePassword)

	// namespace route
	rt.Add("/namespace/getList", http.MethodGet, controller.Namespace{}.GetList)
	rt.Add("/namespace/getTotal", http.MethodGet, controller.Namespace{}.GetTotal)
	rt.Add("/namespace/getBindUserList", http.MethodGet, controller.Namespace{}.GetBindUserList)
	rt.Add("/namespace/getUserOption", http.MethodGet, controller.Namespace{}.GetUserOption)
	rt.Add("/namespace/add", http.MethodPost, controller.Namespace{}.Add).Role(core.RoleAdmin)
	rt.Add("/namespace/edit", http.MethodPost, controller.Namespace{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/addUser", http.MethodPost, controller.Namespace{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/removeUser", http.MethodDelete, controller.Namespace{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager})

	// project route
	rt.Add("/project/getList", http.MethodGet, controller.Project{}.GetList)
	rt.Add("/project/getTotal", http.MethodGet, controller.Project{}.GetTotal)
	rt.Add("/project/getRemoteBranchList", http.MethodGet, controller.Project{}.GetRemoteBranchList)
	rt.Add("/project/getBindServerList", http.MethodGet, controller.Project{}.GetBindServerList)
	rt.Add("/project/getBindUserList", http.MethodGet, controller.Project{}.GetBindUserList)
	rt.Add("/project/getProjectFileList", http.MethodGet, controller.Project{}.GetProjectFileList)
	rt.Add("/project/getProjectFileContent", http.MethodGet, controller.Project{}.GetProjectFileContent)
	rt.Add("/project/add", http.MethodPost, controller.Project{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/edit", http.MethodPost, controller.Project{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/setAutoDeploy", http.MethodPost, controller.Project{}.SetAutoDeploy).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/remove", http.MethodDelete, controller.Project{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/uploadFile", http.MethodPost, controller.Project{}.UploadFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeFile", http.MethodDelete, controller.Project{}.RemoveFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addServer", http.MethodPost, controller.Project{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addUser", http.MethodPost, controller.Project{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeServer", http.MethodDelete, controller.Project{}.RemoveServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeUser", http.MethodDelete, controller.Project{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addFile", http.MethodPost, controller.Project{}.AddFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/editFile", http.MethodPost, controller.Project{}.EditFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addTask", http.MethodPost, controller.Project{}.AddTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeTask", http.MethodPost, controller.Project{}.RemoveTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/getTaskList", http.MethodGet, controller.Project{}.GetTaskList).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/getReviewList", http.MethodGet, controller.Project{}.GetReviewList)

	// monitor route
	rt.Add("/monitor/getList", http.MethodGet, controller.Monitor{}.GetList)
	rt.Add("/monitor/getTotal", http.MethodGet, controller.Monitor{}.GetTotal)
	rt.Add("/monitor/check", http.MethodPost, controller.Monitor{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/add", http.MethodPost, controller.Monitor{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/edit", http.MethodPost, controller.Monitor{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/toggle", http.MethodPost, controller.Monitor{}.Toggle).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/remove", http.MethodDelete, controller.Monitor{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	//// deploy route
	rt.Add("/deploy/getList", http.MethodGet, controller.Deploy{}.GetList)
	rt.Add("/deploy/getPublishTrace", http.MethodGet, controller.Deploy{}.GetPublishTrace)
	rt.Add("/deploy/getPublishTraceDetail", http.MethodGet, controller.Deploy{}.GetPublishTraceDetail)
	rt.Add("/deploy/getCommitList", http.MethodGet, controller.Deploy{}.GetCommitList)
	rt.Add("/deploy/getBranchList", http.MethodGet, controller.Deploy{}.GetBranchList)
	rt.Add("/deploy/getTagList", http.MethodGet, controller.Deploy{}.GetTagList)
	rt.Add("/deploy/getPreview", http.MethodGet, controller.Deploy{}.GetPreview)
	rt.Add("/deploy/review", http.MethodPost, controller.Deploy{}.Review).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/deploy/publish", http.MethodPost, controller.Deploy{}.Publish, middleware.HasPublishAuth)
	rt.Add("/deploy/resetState", http.MethodPost, controller.Deploy{}.ResetState).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/deploy/greyPublish", http.MethodPost, controller.Deploy{}.GreyPublish, middleware.HasPublishAuth).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/deploy/webhook", http.MethodPost, controller.Deploy{}.Webhook, middleware.FilterEvent)
	rt.Add("/deploy/callback", http.MethodGet, controller.Deploy{}.Callback)

	// server route
	rt.Add("/server/getList", http.MethodGet, controller.Server{}.GetList)
	rt.Add("/server/getTotal", http.MethodGet, controller.Server{}.GetTotal)
	rt.Add("/server/getOption", http.MethodGet, controller.Server{}.GetOption)
	rt.Add("/server/check", http.MethodPost, controller.Server{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/add", http.MethodPost, controller.Server{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/edit", http.MethodPost, controller.Server{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/remove", http.MethodDelete, controller.Server{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})

	// crontab route
	rt.Add("/crontab/getList", http.MethodGet, controller.Crontab{}.GetList)
	rt.Add("/crontab/getTotal", http.MethodGet, controller.Crontab{}.GetTotal)
	rt.Add("/crontab/getRemoteServerList", http.MethodGet, controller.Crontab{}.GetRemoteServerList)
	rt.Add("/crontab/getBindServerList", http.MethodGet, controller.Crontab{}.GetBindServerList)
	rt.Add("/crontab/add", http.MethodPost, controller.Crontab{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/edit", http.MethodPost, controller.Crontab{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/import", http.MethodPost, controller.Crontab{}.Import).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/remove", http.MethodDelete, controller.Crontab{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/addServer", http.MethodPost, controller.Crontab{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/crontab/removeCrontabServer", http.MethodDelete, controller.Crontab{}.RemoveCrontabServer).Roles([]string{core.RoleAdmin, core.RoleManager})

	rt.Start()
	return rt
}
