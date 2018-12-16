package route

import (
	"github.com/zhenorzz/goploy/controller"
	router "github.com/zhenorzz/goploy/core"
)

func Init() {
	var rt = new(router.Routes)
	rt.Add("/user/index", controller.UserIndex)
	rt.Add("/github/search", controller.GithubSearch)
	rt.Add("/rsync/add", controller.RsyncAdd)
	rt.Add("/mysql/query", controller.MysqlQuery)
	rt.Start()
}
