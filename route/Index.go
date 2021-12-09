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
	var rt = router.NewRouter()
	// rt.Middleware(example)

	// websocket route
	rt.Add("/ws/connect", http.MethodGet, ws.GetHub().Connect)
	rt.Add("/ws/xterm", http.MethodGet, ws.GetHub().Xterm)
	rt.Add("/ws/sftp", http.MethodGet, ws.GetHub().Sftp)

	// user route
	rt.Add("/user/login", http.MethodPost, controller.User{}.Login).White()
	rt.Add("/user/info", http.MethodGet, controller.User{}.Info)
	rt.Add("/user/getList", http.MethodGet, controller.User{}.GetList)
	rt.Add("/user/getTotal", http.MethodGet, controller.User{}.GetTotal)
	rt.Add("/user/getOption", http.MethodGet, controller.User{}.GetOption)
	rt.Add("/user/add", http.MethodPost, controller.User{}.Add).Role(core.RoleAdmin)
	rt.Add("/user/edit", http.MethodPut, controller.User{}.Edit).Role(core.RoleAdmin)
	rt.Add("/user/remove", http.MethodDelete, controller.User{}.Remove).Role(core.RoleAdmin)
	rt.Add("/user/changePassword", http.MethodPut, controller.User{}.ChangePassword)

	// namespace route
	rt.Add("/namespace/getList", http.MethodGet, controller.Namespace{}.GetList)
	rt.Add("/namespace/getTotal", http.MethodGet, controller.Namespace{}.GetTotal)
	rt.Add("/namespace/getBindUserList", http.MethodGet, controller.Namespace{}.GetBindUserList)
	rt.Add("/namespace/getUserOption", http.MethodGet, controller.Namespace{}.GetUserOption)
	rt.Add("/namespace/add", http.MethodPost, controller.Namespace{}.Add).Role(core.RoleAdmin)
	rt.Add("/namespace/edit", http.MethodPut, controller.Namespace{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/addUser", http.MethodPost, controller.Namespace{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/namespace/removeUser", http.MethodDelete, controller.Namespace{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager})

	// project route
	rt.Add("/project/getList", http.MethodGet, controller.Project{}.GetList)
	rt.Add("/project/getTotal", http.MethodGet, controller.Project{}.GetTotal)
	rt.Add("/project/pingRepos", http.MethodGet, controller.Project{}.PingRepos)
	rt.Add("/project/getRemoteBranchList", http.MethodGet, controller.Project{}.GetRemoteBranchList)
	rt.Add("/project/getBindServerList", http.MethodGet, controller.Project{}.GetBindServerList)
	rt.Add("/project/getBindUserList", http.MethodGet, controller.Project{}.GetBindUserList)
	rt.Add("/project/getProjectFileList", http.MethodGet, controller.Project{}.GetProjectFileList)
	rt.Add("/project/getProjectFileContent", http.MethodGet, controller.Project{}.GetProjectFileContent)
	rt.Add("/project/getReposFileList", http.MethodGet, controller.Project{}.GetReposFileList)
	rt.Add("/project/add", http.MethodPost, controller.Project{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/edit", http.MethodPut, controller.Project{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/setAutoDeploy", http.MethodPut, controller.Project{}.SetAutoDeploy).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/remove", http.MethodDelete, controller.Project{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/uploadFile", http.MethodPost, controller.Project{}.UploadFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeFile", http.MethodDelete, controller.Project{}.RemoveFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addServer", http.MethodPost, controller.Project{}.AddServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addUser", http.MethodPost, controller.Project{}.AddUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeServer", http.MethodDelete, controller.Project{}.RemoveServer).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeUser", http.MethodDelete, controller.Project{}.RemoveUser).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addFile", http.MethodPost, controller.Project{}.AddFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/editFile", http.MethodPut, controller.Project{}.EditFile).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/addTask", http.MethodPost, controller.Project{}.AddTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/removeTask", http.MethodDelete, controller.Project{}.RemoveTask).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/getTaskList", http.MethodGet, controller.Project{}.GetTaskList).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/project/getReviewList", http.MethodGet, controller.Project{}.GetReviewList)

	rt.Add("/repository/getCommitList", http.MethodGet, controller.Repository{}.GetCommitList)
	rt.Add("/repository/getBranchList", http.MethodGet, controller.Repository{}.GetBranchList)
	rt.Add("/repository/getTagList", http.MethodGet, controller.Repository{}.GetTagList)

	// monitor route
	rt.Add("/monitor/getList", http.MethodGet, controller.Monitor{}.GetList)
	rt.Add("/monitor/getTotal", http.MethodGet, controller.Monitor{}.GetTotal)
	rt.Add("/monitor/check", http.MethodPost, controller.Monitor{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/add", http.MethodPost, controller.Monitor{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/edit", http.MethodPut, controller.Monitor{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/toggle", http.MethodPut, controller.Monitor{}.Toggle).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/monitor/remove", http.MethodDelete, controller.Monitor{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	//// deploy route
	rt.Add("/deploy/getList", http.MethodGet, controller.Deploy{}.GetList)
	rt.Add("/deploy/getPublishTrace", http.MethodGet, controller.Deploy{}.GetPublishTrace)
	rt.Add("/deploy/getPublishTraceDetail", http.MethodGet, controller.Deploy{}.GetPublishTraceDetail)
	rt.Add("/deploy/getPreview", http.MethodGet, controller.Deploy{}.GetPreview)
	rt.Add("/deploy/review", http.MethodPut, controller.Deploy{}.Review).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/deploy/resetState", http.MethodPut, controller.Deploy{}.ResetState).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/deploy/publish", http.MethodPost, controller.Deploy{}.Publish, middleware.HasPublishAuth)
	rt.Add("/deploy/rebuild", http.MethodPost, controller.Deploy{}.Rebuild, middleware.HasPublishAuth)
	rt.Add("/deploy/greyPublish", http.MethodPost, controller.Deploy{}.GreyPublish, middleware.HasPublishAuth).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/deploy/webhook", http.MethodPost, controller.Deploy{}.Webhook, middleware.FilterEvent).White()
	rt.Add("/deploy/callback", http.MethodGet, controller.Deploy{}.Callback).White()
	rt.Add("/deploy/fileCompare", http.MethodPost, controller.Deploy{}.FileCompare).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})
	rt.Add("/deploy/fileDiff", http.MethodPost, controller.Deploy{}.FileDiff).Roles([]string{core.RoleAdmin, core.RoleManager, core.RoleGroupManager})

	// server route
	rt.Add("/server/getList", http.MethodGet, controller.Server{}.GetList)
	rt.Add("/server/getTotal", http.MethodGet, controller.Server{}.GetTotal)
	rt.Add("/server/getOption", http.MethodGet, controller.Server{}.GetOption)
	rt.Add("/server/getPublicKey", http.MethodGet, controller.Server{}.GetPublicKey)
	rt.Add("/server/check", http.MethodPost, controller.Server{}.Check).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/add", http.MethodPost, controller.Server{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/edit", http.MethodPut, controller.Server{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/toggle", http.MethodPut, controller.Server{}.Toggle).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/downloadFile", http.MethodGet, controller.Server{}.DownloadFile).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/uploadFile", http.MethodPost, controller.Server{}.UploadFile).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/report", http.MethodGet, controller.Server{}.Report).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/getAllMonitor", http.MethodGet, controller.Server{}.GetAllMonitor)
	rt.Add("/server/addMonitor", http.MethodPost, controller.Server{}.AddMonitor).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/editMonitor", http.MethodPut, controller.Server{}.EditMonitor).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/server/deleteMonitor", http.MethodDelete, controller.Server{}.DeleteMonitor).Roles([]string{core.RoleAdmin, core.RoleManager})

	rt.Add("/cron/report", http.MethodPost, controller.Cron{}.Report).White()
	rt.Add("/cron/getList", http.MethodPost, controller.Cron{}.GetList).White()
	rt.Add("/cron/getLogs", http.MethodPost, controller.Cron{}.GetLogs).White()
	rt.Add("/cron/add", http.MethodPost, controller.Cron{}.Add).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/cron/edit", http.MethodPut, controller.Cron{}.Edit).Roles([]string{core.RoleAdmin, core.RoleManager})
	rt.Add("/cron/remove", http.MethodDelete, controller.Cron{}.Remove).Roles([]string{core.RoleAdmin, core.RoleManager})

	// agent route
	rt.Add("/agent/report", http.MethodPost, controller.Agent{}.Report).White()

	rt.Start()
	return rt
}
