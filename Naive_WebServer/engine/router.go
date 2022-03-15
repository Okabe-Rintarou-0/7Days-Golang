package engine

import (
	"fmt"
	"net/http"
)

type Router struct {
	routerTree *RouterTree
}

func NewRouter() *Router {
	return &Router{
		routerTree: NewRouterTree("/"),
	}
}

func (router *Router) addRoute(pattern string, handler FuncHandler) {
	router.routerTree.Insert(pattern, handler)
}

func (router *Router) handle(c *Context) {
	pattern := c.Method + "-" + c.Path
	if handler, params := router.routerTree.Parse(pattern); handler != nil {
		c.Params = params
		fmt.Printf("Got params: %v\n", params)
		handler(c)
	} else {
		http.Error(c.Writer, "Not Found", http.StatusNotFound)
	}
}
