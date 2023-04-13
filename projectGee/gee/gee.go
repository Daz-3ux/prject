package gee

import (
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
		*RouteGroup						// 继承嵌入类型的所有属性与方法
		router *router
		groups []*RouteGroup 	// store all groups
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

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
