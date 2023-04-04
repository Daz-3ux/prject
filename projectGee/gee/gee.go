package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// Add method to route map (GET-/ = handler)
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

/* GET 与 POST 本质都是 TCP 链接, 因为 HTTP 的规定和 浏览器/服务器 的限制有了区别 */
// GET method: 请求一个指定资源的表示形式，使用 GET 的请求应该只被用于获取数据
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST method: 将实体提交到指定的资源，通常导致在服务器上的状态变化或副作用
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) RUN(addr string)(err error) {
	return http.ListenAndServe(addr, engine)
}

// 在 Go 语言中，实现了接口方法的 struct 都可以强制转换为接口类型
	// handler := (http.Handler)(engine) 											手动转换为接口类型
	// log.Fatal(http.ListenAndServe(":9999", handler))
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		// 设置返回码
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "serverdaz told you : 404 NOT FOUND: %s\n", req.URL)
	}
}