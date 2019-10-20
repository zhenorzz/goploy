package core

import (
	"errors"
	"goploy/model"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Method specifies the HTTP method (GET, POST, PUT, etc.).
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// Goploy callback param
type Goploy struct {
	UserInfo model.User
	Request  *http.Request
	URLQuery url.Values
	Body     []byte
}

// 路由定义
type route struct {
	pattern     string                                          // 正则表达式
	method      string                                          // Method specifies the HTTP method (GET, POST, PUT, etc.).
	roles       []string                                        //允许的角色
	callback    func(w http.ResponseWriter, gp *Goploy)         //Controller函数
	middlewares []func(w http.ResponseWriter, gp *Goploy) error //中间件
}

// Router is route slice and global middlewares
type Router struct {
	Routes      []route
	middlewares []func(w http.ResponseWriter, gp *Goploy) error //中间件
}

// Start a router
func (rt *Router) Start() {
	http.Handle("/", rt)
}

// Add router
// pattern path
// callback  where path should be handle
func (rt *Router) Add(pattern, method string, callback func(w http.ResponseWriter, gp *Goploy), middleware ...func(w http.ResponseWriter, gp *Goploy) error) *Router {
	r := route{pattern: pattern, method: method, callback: callback}
	for _, m := range middleware {
		r.middlewares = append(r.middlewares, m)
	}
	rt.Routes = append(rt.Routes, r)
	return rt
}

// Roles Add many permission to the route
func (rt *Router) Roles(role []string) *Router {
	rt.Routes[len(rt.Routes)-1].roles = append(rt.Routes[len(rt.Routes)-1].roles, role...)
	return rt
}

// Role Add permission to the route
func (rt *Router) Role(role string) *Router {
	rt.Routes[len(rt.Routes)-1].roles = append(rt.Routes[len(rt.Routes)-1].roles, role)
	return rt
}

// Middleware global Middleware handle function
// Example handle parseToken
func (rt *Router) Middleware(middleware func(w http.ResponseWriter, gp *Goploy) error) {
	rt.middlewares = append(rt.middlewares, middleware)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If in production env, serve file in go server,
	// else serve file in npm
	if os.Getenv("ENV") == "production" {
		if "/" == r.URL.Path {
			http.ServeFile(w, r, GlobalPath+"web/dist/index.html")
			return
		}
		files, _ := ioutil.ReadDir(GlobalPath + "web/dist")
		for _, file := range files {
			pattern := "^" + file.Name()
			if match, _ := regexp.MatchString(pattern, r.URL.Path[1:]); match {
				http.ServeFile(w, r, GlobalPath+"web/dist"+r.URL.Path)
				return
			}
		}
	}
	whiteList := map[string]struct{}{
		"/user/login":        {},
		"/user/isShowPhrase": {},
	}
	var userInfo model.User
	if _, ok := whiteList[r.URL.Path]; !ok {
		// check token
		goployTokenCookie, err := r.Cookie(LoginCookieName)
		if err != nil {
			response := Response{Code: IllegalRequest, Message: "Illegal request"}
			response.JSON(w)
			return
		}
		unParseToken := goployTokenCookie.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(unParseToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGN_KEY")), nil
		})

		if err != nil || !token.Valid {
			response := Response{Code: LoginExpired, Message: "登录已过期"}
			response.JSON(w)
			return
		}
		userInfo, err = GetUserInfo(int64(claims["id"].(float64)))
		if err != nil || !token.Valid {
			response := Response{Code: Deny, Message: "获取用户信息失败"}
			response.JSON(w)
			return
		}

		goployTokenStr, err := model.User{ID: int64(claims["id"].(float64)), Name: claims["name"].(string)}.CreateToken()
		if err == nil {
			// update jwt time
			cookie := http.Cookie{Name: LoginCookieName, Value: goployTokenStr, Path: "/", MaxAge: 86400, HttpOnly: true}
			http.SetCookie(w, &cookie)
		}

	}

	// save the body request data because ioutil.ReadAll will clear the requestBody
	var body []byte
	if hasContentType(r, "application/json") {
		body, _ = ioutil.ReadAll(r.Body)
	}
	gp := &Goploy{
		UserInfo: userInfo,
		Request:  r,
		URLQuery: r.URL.Query(),
		Body:     body,
	}
	for _, middleware := range rt.middlewares {
		err := middleware(w, gp)
		if err != nil {
			response := Response{Code: Deny, Message: err.Error()}
			response.JSON(w)
			return
		}
	}
	for _, route := range rt.Routes {
		if route.pattern == r.URL.Path {
			if route.method != r.Method {
				response := Response{Code: Deny, Message: "Invaild request method"}
				response.JSON(w)
				return
			}
			if err := route.hasRole(userInfo.Role); err != nil {
				response := Response{Code: Deny, Message: err.Error()}
				response.JSON(w)
				return
			}
			for _, middleware := range route.middlewares {

				if err := middleware(w, gp); err != nil {
					response := Response{Code: Deny, Message: err.Error()}
					response.JSON(w)
					return
				}
			}

			route.callback(w, gp)
			return
		}
	}

	response := Response{Code: Deny, Message: "No such method"}
	response.JSON(w)
	return

}

func (r *route) hasRole(userRole string) error {
	if len(r.roles) == 0 {
		return nil
	}

	for _, role := range r.roles {
		if role == userRole {
			return nil
		}
	}
	return errors.New("无权限进行此操作")
}

func hasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return false
	}
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
