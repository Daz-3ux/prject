package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 为了代码简洁,为复杂类型起别名
type H map[string]interface{}

// v1:只包含 http.ResponseWriter 和 *http.Request
// v1: 另外提供了对 Method 和 Path 常用属性的直接访问
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index   int 							// 记录当前执行到第几个中间件 
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

// 当在中间件中调用 Next 方法时,控制权交给下一个中间件
// 调用到最后一个中间件后,再从后往前调用每个中间件在 Next 方法之后定义的部分
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ { // 妙不可言
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	//value, _ := c.Params[key]
	value := c.Params[key]
	return value
}

// PostForm 方法: 返回查询中指定组件的第一个值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 方法: 从HTTP请求的URL查询参数中获取指定键名对应的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 方法: 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 方法: 修改上下文中的响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 将一个 字符串 写入响应体
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 将一个 JSON 写入响应体
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 将一个 字节数组 写入响应体
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 将一个 HTML 字符串写入响应体
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
