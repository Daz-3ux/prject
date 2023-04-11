package gee

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type (
	RouteGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouteGroup   // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	// 整个框架的资源都是由 Engine 统一协调的
	Engine struct {
		*RouteGroup
		router *router
		groups []*RouteGroup // store all groups
	}
)

func New() *Engine {
	//return &Engine{router: newRouter()}
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}

	return engine
}

func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.engine
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouteGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// just for debug(with a ugly way)
func (engine *Engine) returnRoute(method string, pattern string) {
	node, params := engine.router.getRoute(method, pattern)
	fmt.Println("node = ", node, "params = ", params)
}

// 接管所有的 HTTP 请求
// Go 语言中: HTTP 服务器可以通过实现 http.Handler 接口来处理 HTTP 请求,
// 这个接口只有一个方法 ServeHTTP -- 用来处理请求并向客户端发送响应
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 填充上下文
	c := newContext(w, req)
	// 启动处理部分
	engine.router.handle(c)
}
