package engine

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	roots map[string]*Node
}

func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*Node),
	}
}

func ParsePattern(pattern string) []string {
	tokens := strings.Split(pattern, "/")
	for i := 0; i < len(tokens); {
		token := tokens[i]
		if len(token) == 0 {
			tokens = append(tokens[:i], tokens[i+1:]...)
		} else {
			i++
		}
	}
	fmt.Printf("got tokens[size = %d]: %v\n", len(tokens), tokens)
	return tokens
}

func (router *Router) addRoute(method string, pattern string, handler FuncHandler) {
	root, ok := router.roots[method]
	if !ok {
		root = NewNode(method, 0)
		router.roots[method] = root
	}
	root.Insert(ParsePattern(pattern), handler)
}

func (router *Router) handle(c *Context) {
	root, ok := router.roots[c.Method]
	if !ok {
		return
	}

	for _, middleware := range c.Middlewares {
		middleware(c)
	}

	if handler, params := root.Parse(ParsePattern(c.Pattern)); handler != nil {
		c.Params = params
		fmt.Printf("Got params: %v\n", params)
		handler(c)
	} else {
		http.Error(c.Writer, "Not Found", http.StatusNotFound)
	}
}
