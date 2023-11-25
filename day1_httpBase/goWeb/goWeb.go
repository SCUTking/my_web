package goWeb

import (
	"fmt"
	"net/http"
)

// Engine is the uni handler for all requests
const separator string = "-"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc, 0)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	route := method + separator + pattern
	engine.router[route] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(add string) (err error) {
	//engine实现了handler接口，利用engine接收全部的http请求
	return http.ListenAndServe(add, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//根据是否请求的方法和路径组装key  ————静态的
	key := req.Method + separator + req.URL.Path

	//找到对应的handler处理
	handlerFunc, ok := engine.router[key]
	if ok {
		handlerFunc(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
