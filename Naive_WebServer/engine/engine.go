package engine

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type FuncHandler func(c *Context)

type Engine struct {
	router *Router
}

func NewEngine() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

func (engine *Engine) addRoute(pattern string, handler FuncHandler) {
	engine.router.addRoute(pattern, handler)
	fmt.Printf("Add a router[pattern = %s] to the engine\n", pattern)
}

func (engine *Engine) Bind(method string, path string, handler FuncHandler) {
	engine.addRoute(method+"-"+path, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	engine.router.handle(NewContext(w, r))
}

func (engine *Engine) Run(port uint16) {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), engine))
}
