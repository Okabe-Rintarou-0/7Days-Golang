package engine

import (
	"fmt"
	"net/http"
	"path"
)

type Group struct {
	prefix       string
	parent       *Group
	engine       *Engine
	middlewares  []FuncHandler
	interceptors []Interceptor
}

func NewGroup(prefix string, parent *Group, engine *Engine) *Group {
	return &Group{
		prefix: prefix,
		parent: parent,
		engine: engine,
	}
}

func (g *Group) Group(prefix string) *Group {
	newGroup := g.engine.Group(prefix)
	newGroup.parent = g
	return newGroup
}

func (g *Group) AddMiddleWare(middleware FuncHandler) *Group {
	g.middlewares = append(g.middlewares, middleware)
	return g
}

func (g *Group) AddMiddleWares(middlewares []FuncHandler) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

func (g *Group) AddInterceptor(interceptor Interceptor) *Group {
	g.interceptors = append(g.interceptors, interceptor)
	return g
}

func (g *Group) AddInterceptors(interceptors []Interceptor) *Group {
	g.interceptors = append(g.interceptors, interceptors...)
	return g
}

func (g *Group) Prefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.Prefix() + g.prefix
}

func (g *Group) addRoute(method, pattern string, handler FuncHandler) {
	pattern = g.Prefix() + pattern
	g.engine.addRoute(method, pattern, handler)
}

func (g *Group) Get(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodGet, pattern, handler)
	return g
}

func (g *Group) Post(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodPost, pattern, handler)
	return g
}

func (g *Group) Delete(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodHead, pattern, handler)
	return g
}

func (g *Group) Head(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodHead, pattern, handler)
	return g
}

func (g *Group) Put(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodPut, pattern, handler)
	return g
}

func (g *Group) Options(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodOptions, pattern, handler)
	return g
}

func (g *Group) Connect(pattern string, handler FuncHandler) *Group {
	g.addRoute(http.MethodConnect, pattern, handler)
	return g
}

func (g *Group) getStaticHandler(relativePath string, fs http.FileSystem) FuncHandler {
	absolutePath := path.Join(g.Prefix(), relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		filePath := c.Param("filePath")
		if _, err := fs.Open(filePath); err != nil {
			fmt.Println(err.Error())
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (g *Group) BindStatic(relativePath string, root string) *Group {
	handler := g.getStaticHandler(relativePath, http.Dir(root))
	return g.Get(path.Join(relativePath, "/*filePath"), handler)
}
