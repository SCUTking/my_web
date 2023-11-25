package goWeb

import (
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
const separator string = "-"

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
	router *router
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newRouteGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newRouteGroup)
	return newRouteGroup
}

func (group *RouterGroup) addRoute(method string, com string, handler HandlerFunc) {
	pattern := group.prefix + com
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (engine *Engine) Run(add string) (err error) {
	//engine实现了handler接口，利用engine接收全部的http请求
	return http.ListenAndServe(add, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	c := newContext(w, req)
	engine.router.handle(c)
}
