package core

import (
	"net/http"
)

// 路由定义
type route struct {
	pattern  string                                       // 正则表达式
	callback func(w http.ResponseWriter, r *http.Request) //Controller函数
}

type Routes []route

func (rt *Routes) router(w http.ResponseWriter, r *http.Request) {
	for _, route := range *rt {
		if route.pattern == r.URL.Path {
			route.callback(w, r)
		}
	}
}

func (rt *Routes) Add(pattern string, callback func(w http.ResponseWriter, r *http.Request)) {
	*rt = append(*rt, route{pattern, callback})
}

func (rt *Routes) Start() {
	http.HandleFunc("/", rt.router)
}
