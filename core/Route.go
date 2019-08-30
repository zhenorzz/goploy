package core

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenInfo pasre the jwt
type TokenInfo struct {
	ID   int64
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
	roles       []string                                        //允许的角色
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

// Roles Add many permision to the route
func (rt *Router) Roles(role []string) *Router {
	rt.Routes[len(rt.Routes)-1].roles = append(rt.Routes[len(rt.Routes)-1].roles, role...)
	return rt
}

// Role Add permision to the route
func (rt *Router) Role(role string) *Router {
	rt.Routes[len(rt.Routes)-1].roles = append(rt.Routes[len(rt.Routes)-1].roles, role)
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
			ID:   int64(claims["id"].(float64)),
			Name: claims["name"].(string),
		}
	}

	// save the body request data bucause ioutil.ReadAll will clear the requestBody
	var body []byte
	if hasContentType(r, "application/json") {
		body, _ = ioutil.ReadAll(r.Body)
	}
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
			if err := route.hasRole(tokenInfo.ID); err != nil {
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

func (r *route) hasRole(userID int64) error {
	if len(r.roles) == 0 {
		return nil
	}
	userInfo, err := GetUserInfo(userID)
	if err != nil {
		return err
	}

	for _, role := range r.roles {
		if role == userInfo.Role {
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
