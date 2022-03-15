package engine

import (
	"strings"
)

type RouterTree struct {
	root *Node
}

func NewRouterTree(token string) *RouterTree {
	return &RouterTree{
		root: NewNode(token, 0),
	}
}

func (rt *RouterTree) Insert(pattern string, handler FuncHandler) {
	rt.root.Insert(strings.Split(pattern, "/"), handler)
}

func (rt *RouterTree) Parse(pattern string) (FuncHandler, map[string]string) {
	return rt.root.Parse(strings.Split(pattern, "/"))
}
