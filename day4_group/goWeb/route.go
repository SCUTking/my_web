package goWeb

import (
	"net/http"
	"strings"
)

type router struct {
	route    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc, 0),
		route:    make(map[string]*node, 0),
	}
}
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + separator + pattern

	parts := parsePattern(pattern)
	_, ok := router.route[method]
	if !ok {
		router.route[method] = &node{}
	}

	router.route[method].insert(pattern, parts, 0)

	router.handlers[key] = handler

}

// 有关键字：的进行提取
func (router *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	//找到对应方法的根节点
	root, ok := router.route[method]
	if !ok {
		return nil, nil
	}

	//找到对应的路径
	n := root.search(searchParts, 0)

	params := make(map[string]string)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//解析 路径中的冒号部分     换成实际的部分
			//例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			//
			//例如  /static/*filepath，可以匹配/static/js/jQuery.js
			//要将  js jQuery 部分重新组装成  /js/jQuery.js
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}
	return nil, nil

}

func (router *router) handle(c *Context) {
	n, params := router.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		//key := c.Method + "-" + c.Path    不能使用这个  因为动态的有些比对不了
		//比如  /a 实际上是/*  直接使用  /a 没有这个东西
		key := c.Method + "-" + n.pattern
		router.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

// 将路径中的/给去除掉，使用一个切片将“/”之间的路径字符串记录下来
// "/p/:lang"-->"{p :lang}"
// "/p/*file/a"-->"p *file"
func parsePattern(pattern string) []string {
	//根据/分割
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//遇到第一个*号就结束了，一个路径只能有一个*
			if item[0] == '*' {
				break
			}
		}
	}

	return parts

}
