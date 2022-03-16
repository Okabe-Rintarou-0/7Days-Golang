package engine

import (
	"fmt"
	"strings"
)

type Node struct {
	Token    string
	Children []*Node
	WildCard uint8
	handler  FuncHandler
}

func NewNode(token string, wildcard uint8) *Node {
	return &Node{
		Token:    token,
		Children: []*Node{},
		WildCard: wildcard,
		handler:  nil,
	}
}

func (root *Node) Info() string {
	info := fmt.Sprintf("[Token] = %s\n[Child num] = %d\n", root.Token, len(root.Children))
	for i, child := range root.Children {
		info += fmt.Sprintf("[Child %d's token] = %s\n", i, child.Token)
	}
	return info
}

func (root *Node) Insert(tokens []string, handler FuncHandler) {
	p := root
	var nextChild *Node
	var found bool
	var childToken string
	for _, curToken := range tokens {
		found = false
		for _, child := range p.Children {
			childToken = child.Token
			if childToken == curToken || string(rune(child.WildCard))+childToken == curToken {
				nextChild = child
				found = true
				break
			}
		}

		if !found {
			var wildCard uint8 = 0

			if strings.HasPrefix(curToken, ":") || strings.HasPrefix(curToken, "*") {
				wildCard = curToken[0]
				curToken = curToken[1:]
			}
			nextChild = NewNode(curToken, wildCard)
			p.Children = append(p.Children, nextChild)
			if wildCard == '*' {
				p = nextChild
				break
			}
		}
		p = nextChild
	}

	p.handler = handler
}

func (root *Node) Parse(tokens []string) (FuncHandler, map[string]string) {
	p := root
	params := make(map[string]string)
	var i int
	for t, curToken := range tokens {
		numChildren := len(p.Children)
		for i = 0; i < numChildren; i++ {
			child := p.Children[i]
			childToken := child.Token
			if childToken == curToken {
				p = child
				break
			} else if child.WildCard > 0 {
				if child.WildCard == '*' {
					params[childToken] = strings.Join(tokens[t:], "/")
					return child.handler, params
				}
				params[childToken] = curToken
				p = child
				break
			}
		}
		if i == numChildren {
			return nil, nil
		}
	}
	//fmt.Println(p.Info())
	return p.handler, params
}
