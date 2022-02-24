package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"github.com/zhenorzz/goploy/web"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Goploy callback param
type Goploy struct {
	UserInfo       model.User
	Namespace      model.Namespace
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	URLQuery       url.Values
	Body           []byte
}

type RouteApi interface {
	Routes() []Route
}

type Response interface {
	Write(http.ResponseWriter) error
}

type Route struct {
	pattern     string
	method      string                    // Method specifies the HTTP method (GET, POST, PUT, etc.).
	roles       map[string]struct{}       // permission role
	white       bool                      // no need to login
	middlewares []func(gp *Goploy) error  // Middlewares run before callback, trigger error will end the request
	callback    func(gp *Goploy) Response // Controller function
	logFunc     func(gp *Goploy, resp Response)
}

// Router is Route slice and global middlewares
type Router struct {
	routes      map[string]Route
	middlewares []func(gp *Goploy) error // Middlewares run before all Route
}

func NewRouter() Router {
	return Router{
		routes: map[string]Route{},
	}
}

func NewRoute(pattern, method string, callback func(gp *Goploy) Response) Route {
	return Route{
		pattern:  pattern,
		method:   method,
		callback: callback,
		roles:    map[string]struct{}{},
	}
}

func NewWhiteRoute(pattern, method string, callback func(gp *Goploy) Response) Route {
	return Route{
		pattern:  pattern,
		method:   method,
		white:    true,
		roles:    map[string]struct{}{},
		callback: callback,
	}
}

// Start a router
func (rt Router) Start() {
	if config.Toml.Env == "production" {
		subFS, err := fs.Sub(web.Dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/assets/", http.FileServer(http.FS(subFS)))
		http.Handle("/favicon.ico", http.FileServer(http.FS(subFS)))
	}
	http.Handle("/", rt)
}

// Middleware global Middleware handle function
func (rt Router) Middleware(middleware func(gp *Goploy) error) {
	rt.middlewares = append(rt.middlewares, middleware)
}

// Add pattern path
// callback where path should be handled
func (rt Router) Add(ra RouteApi) Router {
	for _, r := range ra.Routes() {
		rt.routes[r.pattern] = r
	}
	return rt
}

// Roles Add much permission to the Route
func (r Route) Roles(roles ...string) Route {
	for _, role := range roles {
		r.roles[role] = struct{}{}
	}
	return r
}

// Middleware global Middleware handle function
func (r Route) Middleware(middleware func(gp *Goploy) error) Route {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// LogFunc callback finished
func (r Route) LogFunc(f func(gp *Goploy, resp Response)) Route {
	r.logFunc = f
	return r
}

func (rt Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If in production env, serve file in go server,
	// else serve file in npm
	if config.Toml.Env == "production" {
		if "/" == r.URL.Path {
			r, err := web.Dist.Open("dist/index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer r.Close()
			contents, err := ioutil.ReadAll(r)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, string(contents))
			return
		}
	}

	_, resp := rt.doRequest(w, r)
	if err := resp.Write(w); err != nil {
		Log(ERROR, err.Error())
	}
	return
}

func (rt Router) doRequest(w http.ResponseWriter, r *http.Request) (*Goploy, Response) {
	gp := new(Goploy)
	route, ok := rt.routes[r.URL.Path]
	if !ok {
		return gp, response.JSON{Code: response.Deny, Message: "No such method"}
	}
	if route.method != r.Method {
		return gp, response.JSON{Code: response.IllegalRequest, Message: "Invalid request method"}
	}

	if !route.white {
		unParseToken := ""
		// check token
		goployTokenCookie, err := r.Cookie(config.Toml.Cookie.Name)
		if err != nil {
			unParseToken = r.URL.Query().Get(config.Toml.Cookie.Name)
		} else {
			unParseToken = goployTokenCookie.Value
		}

		if unParseToken == "" {
			return gp, response.JSON{Code: response.IllegalRequest, Message: "Illegal request"}
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(unParseToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Toml.JWT.Key), nil
		})

		if err != nil || !token.Valid {
			return gp, response.JSON{Code: response.LoginExpired, Message: "Login expired"}
		}

		namespaceIDRaw := r.Header.Get(NamespaceHeaderName)
		if namespaceIDRaw == "" {
			namespaceIDRaw = r.URL.Query().Get(NamespaceHeaderName)
		}

		namespaceID, err := strconv.ParseInt(namespaceIDRaw, 10, 64)
		if err != nil {
			return gp, response.JSON{Code: response.Deny, Message: "Invalid namespace"}
		}

		gp.Namespace, err = model.Namespace{
			ID:     namespaceID,
			UserID: int64(claims["id"].(float64)),
		}.GetDataByUserNamespace()

		if err != nil {
			if err == sql.ErrNoRows {
				return gp, response.JSON{Code: response.NamespaceInvalid, Message: "No available namespace"}
			} else {
				return gp, response.JSON{Code: response.Deny, Message: err.Error()}
			}
		}

		if err = route.hasRole(gp.Namespace.Role); err != nil {
			return gp, response.JSON{Code: response.Deny, Message: err.Error()}
		}

		gp.UserInfo, err = model.User{ID: int64(claims["id"].(float64))}.GetData()
		if err != nil {
			return gp, response.JSON{Code: response.Deny, Message: "Get user information error"}
		}

		goployTokenStr, err := model.User{ID: int64(claims["id"].(float64)), Name: claims["name"].(string)}.CreateToken()
		if err == nil {
			// update jwt time
			cookie := http.Cookie{Name: config.Toml.Cookie.Name, Value: goployTokenStr, Path: "/", MaxAge: config.Toml.Cookie.Expire, HttpOnly: true}
			http.SetCookie(w, &cookie)
		}

	}

	gp.Request = r
	gp.ResponseWriter = w
	gp.URLQuery = r.URL.Query()

	// save the body request data because ioutil.ReadAll will clear the requestBody
	if r.ContentLength > 0 && hasContentType(r, "application/json") {
		gp.Body, _ = ioutil.ReadAll(r.Body)
	}

	// common middlewares
	for _, middleware := range rt.middlewares {
		err := middleware(gp)
		if err != nil {
			return gp, response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	// route middlewares
	for _, middleware := range route.middlewares {
		if err := middleware(gp); err != nil {
			return gp, response.JSON{Code: response.Error, Message: err.Error()}
		}
	}

	resp := route.callback(gp)

	if route.logFunc != nil {
		route.logFunc(gp, resp)
	}

	return gp, resp
}

func (r Route) hasRole(namespaceRole string) error {
	if len(r.roles) == 0 {
		return nil
	}

	if _, ok := r.roles[namespaceRole]; ok {
		return nil
	}
	return errors.New("no permission")
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
