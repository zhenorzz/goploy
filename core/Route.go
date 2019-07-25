package core

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhenorzz/goploy/model"
)

// TokenInfo pasre the jwt
type TokenInfo struct {
	ID   uint32
	Name string
}

// Goploy callback param
type Goploy struct {
	TokenInfo TokenInfo
	Request   *http.Request
	URLQuery  url.Values
	Body      []byte
}

// 路由定义
type route struct {
	pattern     string                                          // 正则表达式
	auth        []string                                        //权限
	callback    func(w http.ResponseWriter, gp *Goploy)         //Controller函数
	middlewares []func(w http.ResponseWriter, gp *Goploy) error //中间件
}

// Router is route slice and golbal middlewares
type Router struct {
	Routes      []route
	middlewares []func(w http.ResponseWriter, gp *Goploy) error //中间件
}

// Start a router
func (rt *Router) Start() {
	http.HandleFunc("/", rt.router)
}

// Add router
// pattern path
// callback  where path should be handle
func (rt *Router) Add(pattern string, callback func(w http.ResponseWriter, gp *Goploy), middleware ...func(w http.ResponseWriter, gp *Goploy) error) *Router {
	r := route{pattern: pattern, callback: callback}
	for _, m := range middleware {
		r.middlewares = append(r.middlewares, m)
	}
	rt.Routes = append(rt.Routes, r)
	return rt
}

// AuthMany Add many permision to the route
func (rt *Router) AuthMany(auth []string) *Router {
	rt.Routes[len(rt.Routes)-1].auth = append(rt.Routes[len(rt.Routes)-1].auth, auth...)
	return rt
}

// Auth Add permision to the route
func (rt *Router) Auth(auth string) *Router {
	rt.Routes[len(rt.Routes)-1].auth = append(rt.Routes[len(rt.Routes)-1].auth, auth)
	return rt
}

// Middleware golbal Middleware handle function
// Example handle praseToken
func (rt *Router) Middleware(middleware func(w http.ResponseWriter, gp *Goploy) error) {
	rt.middlewares = append(rt.middlewares, middleware)
}

func (rt *Router) router(w http.ResponseWriter, r *http.Request) {
	// If in production env, serve file in go server,
	// else serve file in npm
	if os.Getenv("ENV") == "production" {
		if "/" == r.URL.Path {
			http.ServeFile(w, r, GolbalPath+"web/dist/index.html")
			return
		}
		files, _ := ioutil.ReadDir(GolbalPath + "web/dist")
		for _, file := range files {
			pattern := "^" + file.Name()
			if match, _ := regexp.MatchString(pattern, r.URL.Path[1:]); match {
				http.ServeFile(w, r, GolbalPath+"web/dist"+r.URL.Path)
				return
			}
		}
	}
	whiteList := map[string]struct{}{
		"/user/login":        {},
		"/user/isShowPhrase": {},
	}
	var tokenInfo TokenInfo
	if _, ok := whiteList[r.URL.Path]; !ok {
		// check token
		goployTokenCookie, err := r.Cookie("goploy_token")
		if err != nil {
			response := Response{Code: 10001, Message: "非法请求"}
			response.JSON(w)
			return
		}
		unPraseToken := goployTokenCookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(unPraseToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGN_KEY")), nil
		})

		if err != nil || !token.Valid {
			response := Response{Code: 10086, Message: "登录已过期"}
			response.JSON(w)
			return
		}
		tokenInfo = TokenInfo{
			ID:   uint32(claims["id"].(float64)),
			Name: claims["name"].(string),
		}
	}

	// save the body request data bucause ioutil.ReadAll will clear the requestBody
	body, _ := ioutil.ReadAll(r.Body)
	gp := &Goploy{
		TokenInfo: tokenInfo,
		Request:   r,
		URLQuery:  r.URL.Query(),
		Body:      body,
	}
	for _, middleware := range rt.middlewares {
		err := middleware(w, gp)
		if err != nil {
			response := Response{Code: 1, Message: err.Error()}
			response.JSON(w)
			return
		}
	}
	for _, route := range rt.Routes {
		if route.pattern == r.URL.Path {
			if err := route.hasPermission(tokenInfo.ID); err != nil {
				response := Response{Code: 1, Message: err.Error()}
				response.JSON(w)
				return
			}
			for _, middleware := range route.middlewares {

				if err := middleware(w, gp); err != nil {
					response := Response{Code: 1, Message: err.Error()}
					response.JSON(w)
					return
				}
			}

			route.callback(w, gp)
		}
	}
}

func (r *route) hasPermission(userID uint32) error {
	if len(r.auth) == 0 {
		return nil
	}
	var permissions model.Permissions

	if x, found := Cache.Get("permissions:" + strconv.Itoa(int(userID))); found {
		permissions = *x.(*model.Permissions)
	} else {
		return errors.New("token过期")
	}
	for _, permission := range permissions {
		for _, auth := range r.auth {
			if permission.URI == auth {
				return nil
			}
		}
	}

	return errors.New("无权限进行此操作")
}
