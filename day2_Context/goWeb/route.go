package goWeb

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc, 0)}
}
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	route := method + separator + pattern
	router.handlers[route] = handler
}

func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
