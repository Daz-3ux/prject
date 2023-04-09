package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
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
