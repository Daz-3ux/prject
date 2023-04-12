package main

import (
	//"fmt"
	"net/http"

	"gee"
)

func main() {
	r := gee.New()

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v1 := r.Group("/v1") 
	{ // just for read
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=daz
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{	// just for read
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/daz
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	//fmt.Printf("%p -> %v, %p -> %v, %p -> %v\n", r.Groups[0], *r.Groups[0], r.Groups[1], *r.Groups[1], r.Groups[2], *r.Groups[2])

	// RUN start webserver
	r.RUN(":9999")
}
