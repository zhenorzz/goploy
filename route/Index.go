package route

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/controller"
)

func Init()  {
	var rt = new(router.Routes)
	rt.Add("/user/index", controller.UserIndex)
	rt.Add("/github/search", controller.GithubSearch)
	rt.Add("/rsync/add", controller.RsyncAdd)
	rt.Start()
}