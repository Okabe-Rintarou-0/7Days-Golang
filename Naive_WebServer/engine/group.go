package engine

type Group struct {
	prefix      string
	parent      *Group
	engine      *Engine
	middlewares []FuncHandler
}

func NewGroup(prefix string, parent *Group, engine *Engine) *Group {
	return &Group{
		prefix: prefix,
		parent: parent,
		engine: engine,
	}
}

func (g *Group) NewGroup(prefix string) *Group {
	newGroup := g.engine.NewGroup(prefix)
	newGroup.parent = g
	return newGroup
}

func (g *Group) AddMiddleWare(middleware FuncHandler) *Group {
	g.middlewares = append(g.middlewares, middleware)
	return g
}

func (g *Group) AddMiddleWares(middlewares []FuncHandler) *Group {
	for _, middleware := range middlewares {
		g.middlewares = append(g.middlewares, middleware)
	}
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

func (g *Group) Bind(method string, pattern string, handler FuncHandler) *Group {
	g.addRoute(method, pattern, handler)
	return g
}
