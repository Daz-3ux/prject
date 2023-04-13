package gee

import (
	//"fmt"
	//"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

// 构建支持分组和中间件的 HTTP 框架
type (
	// 加入分组: 以相同前缀区分
	// 中间件支持: 中间件可以给框架提供无限扩展能力(与分组绑定,支持中间件链式调用)
	RouteGroup struct { // 代表一个路由分组
		// 该分组的路由前缀
		prefix      string
		// 该分组的中间件列表,支持中间件链式调用
		middlewares []HandlerFunc // support middleware
		// 该分组的父级分组
		parent      *RouteGroup   // support nesting
		// 整个框架的资源都是由 Engine 统一协调的,可以通过 Engine 间接的访问各类接口
		engine      *Engine       // all groups share a Engine instance
	}

	// 最顶层分组: 整个框架的资源都是由 Engine 统一协调的
	Engine struct {
		// 用于管理顶层的路由分组
		*RouteGroup						// 继承了 RouteGroup 的所有属性和方法(结构体的继承)
		// 处理 HTTP 请求
		router *router				// 路由器实例
		// 用于存储管理所有的路由分组
		Groups []*RouteGroup 	// 所有分组的列表
	}
)

/*
	r := gee.New()
	v1 := r.Group("/v1")
*/

func New() *Engine {
	// 初始化一个新的 Engine 实例
	engine := &Engine{router: newRouter()}
	// 创建一个新的路由实例
	engine.RouteGroup = &RouteGroup{engine: engine}
	// 创建一个顶层路由分组
	engine.Groups = []*RouteGroup{engine.RouteGroup} // Groups中的第一个元素在此处添加,为 nil

	return engine
}

// 创建一个新的路由分组,并将其添加到路由器实例的分组列表中
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.engine
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	//fmt.Printf("NewGroup: prefix=%s, parent=%v, engine=%v\n", newGroup.prefix, newGroup.parent, newGroup.engine)
	engine.Groups = append(engine.Groups, newGroup)
	return newGroup
}

func (group *RouteGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 将中间件应用到某个 Group
func (group *RouteGroup) Use(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

func (engine *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// // just for debug(with a ugly way)
// func (engine *Engine) returnRoute(method string, pattern string) {
// 	node, params := engine.router.getRoute(method, pattern)
// 	fmt.Println("node = ", node, "params = ", params)
// }

// 接管所有的 HTTP 请求
// Go 语言中: HTTP 服务器可以通过实现 http.Handler 接口来处理 HTTP 请求,
// 这个接口只有一个方法 ServeHTTP -- 用来处理请求并向客户端发送响应
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.Groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 填充上下文
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
