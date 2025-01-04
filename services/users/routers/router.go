package routers

import (
	"net/http"
	"strings"
)

type Router struct {
  BasePath string
	Routes map[string]Route
}

type Route struct {
  Middleware string // later a function
  Handler http.HandlerFunc
}

func newRouter(basepath string) Router {
  return Router{
    BasePath: basepath,
    Routes: make(map[string]Route),
  }
}

func (r *Router) addRoute(path string,method string, handler http.HandlerFunc) {
  r.Routes[path+":"+method] = Route{
    Handler: handler,
  }
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  if !strings.HasPrefix(req.URL.Path, r.BasePath) {
    http.NotFound(w,req)
    return
  }
  path := strings.Split(req.URL.Path,"/")
  next_path := path[len(path)-1]+":"+req.Method
  route,exists := r.Routes[next_path]
  if !exists {
    http.NotFound(w,req)
    return
  }
  route.Handler(w,req)
}
