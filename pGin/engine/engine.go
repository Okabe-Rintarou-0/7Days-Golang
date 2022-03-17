package engine

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FuncHandler func(c *Context)
type Interceptor func(c *Context) bool

type Engine struct {
	groups        []*Group
	router        *Router
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{
		router: NewRouter(),
	}
	engine.groups = append(engine.groups, NewGroup("/", nil, engine))
	return engine
}

func (engine *Engine) Group(prefix string) *Group {
	newGroup := NewGroup(prefix, nil, engine)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *Engine) Default() *Group {
	return engine.groups[0]
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
	c.engine = engine
	for _, group := range engine.groups {
		if strings.HasPrefix(c.Pattern, group.Prefix()) {
			c.Middlewares = append(c.Middlewares, group.middlewares...)
			c.Interceptors = append(c.Interceptors, group.interceptors...)
		}
	}

	defer recoverPanic()

	engine.router.handle(c)
}

func (engine *Engine) Run(port uint16) {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), engine))
}

func (engine *Engine) LoadHTMLGlob(path string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(path))
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}
