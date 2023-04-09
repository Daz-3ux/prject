package main

import (
	"net/http"

	"gee"
	"runtime/debug"
)

func main() {
	debug.SetTraceback("all")
	r := gee.New();

	// hanler 的参数变为 gee.Context
	// 提供了查询 Query/PostForm 参数功能
	// 封装了 HTML/String/JSON, 能够快速构造 HTTP 响应
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=daz
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		// expect /hello/daz
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	// RUN start webserver
	r.RUN(":9999")
}
