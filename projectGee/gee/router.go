package gee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type router struct {
	// 使用 roots 来存储每种请求的 Tire树 根节点
	// eg: roots['SET'] roots['POST']
	roots map[string]*node
	// 使用 handlers 存储每种请求的 HandlerFunc
	// eg: handlers['GET=/p/:lang/doc'] handlers['POST-/p/book']
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// 加路由方法
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	// parts
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	fmt.Println("method: ", method, "parts: ", parts, "key: ", key)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 用于获取指定请求方法和路径对应的路由节点和参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 解析请求路径,获取参数路径
	searchPatrs := parsePattern(path)
	params := make(map[string]string)

	// 获取指定请求方法的根节点
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 在根节点下搜索指定路径
	n := root.search(searchPatrs, 0)

	if n != nil {
		// 解析路由模式,获取参数名称和通配符
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 如果路由模式以冒号开头,则改部分是一个参数
			// /users/:id，URL路径是 /users/123: 则参数名称是 "id"，参数值是 "123"
			if part[0] == ':' {
				params[part[1:]] = searchPatrs[index]
			}
			// 如果路由模式以星号开头,并且长度长度大于1,则该部分是匹配剩余路径的通配符
			// /users/*path，URL路径是 /users/123/profile: 参数名称是 "path"，参数值是 "123/profile"
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPatrs[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	// 在路由中获取 路由节点 和 参数
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		fmt.Println("HANDLE: params: ", params, "key: ", key, "method: ", c.Method, "pattern: ", n.pattern)
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "serverdaz told you: 404 NOT FOUND: %s\n", c.Path)
	}
}
