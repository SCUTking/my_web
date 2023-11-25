package my_web

import (
	"fmt"
	goWeb2 "my_web/day1_httpBase/goWeb"
	"my_web/day2_Context/goWeb"
	"net/http"
	"testing"
)

func TestGoWeb(t *testing.T) {
	r := goWeb2.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}

func TestContext(t *testing.T) {
	r := goWeb.New()
	r.GET("/", func(c *goWeb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *goWeb.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *goWeb.Context) {
		c.JSON(http.StatusOK, goWeb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
