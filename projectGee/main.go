package main

import (
	"fmt"
	"net/http"

	"gee"
)

func main() {
	// New create gee instance
	r := gee.New();

	// GET add route(now only support static route)
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	// RUN start webserver
	r.RUN(":9999")
}