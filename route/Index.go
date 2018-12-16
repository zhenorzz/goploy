package route

import (
	"github.com/zhenorzz/goploy/controller"
	router "github.com/zhenorzz/goploy/core"
)

func Init() {
	var rt = new(router.Routes)
	rt.Add("/user/index", new(controller.User).Index)
	rt.Add("/github/search", new(controller.Github).Search)
	rt.Add("/rsync/add", new(controller.Rsync).Add)
	rt.Add("/mysql/query", new(controller.Mysql).Query)
	rt.Start()
}
