package engine

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FuncHandler func(c *Context)

type Engine struct {
	groups []*Group
	router *Router
}

func NewEngine() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

func (engine *Engine) NewGroup(prefix string) *Group {
	newGroup := NewGroup(prefix, nil, engine)
	engine.groups = append(engine.groups)
	return newGroup
}

func (engine *Engine) addRoute(method, pattern string, handler FuncHandler) {
	fmt.Printf("Add a router[pattern = %s] to the engine\n", pattern)
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) bind(method string, pattern string, handler FuncHandler) {
	engine.addRoute(method, pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)

	for _, group := range engine.groups {
		if strings.HasPrefix(c.Pattern, group.Prefix()) {
			c.Middlewares = append(c.Middlewares, group.middlewares...)
		}
	}

	engine.router.handle(c)
}

func (engine *Engine) Run(port uint16) {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), engine))
}
