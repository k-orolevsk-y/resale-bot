package bot

import (
	"regexp"
)

type methodTrees []*methodTree

func (trees methodTrees) get(method string) (*methodTree, bool) {
	for _, tree := range trees {
		if tree.method == method {
			return tree, true
		}
	}

	return nil, false
}

func (trees methodTrees) compileHandlersChain(method, query string) HandlersChain {
	var handlers HandlersChain
	onlyMiddlewares := true

	for _, tree := range trees {
		if tree.method == method || tree.method == "*" {
			for _, n := range tree.nodes {
				if n.regex.MatchString(query) {
					if tree.method != "*" {
						onlyMiddlewares = false
					}

					handlers = append(handlers, n.handlers...)
				}
			}
		}
	}

	if onlyMiddlewares {
		return HandlersChain{}
	}

	return handlers
}

type methodTree struct {
	method string
	nodes  []*node
}

func (tree *methodTree) get(regex *regexp.Regexp) (*node, bool) {
	for _, n := range tree.nodes {
		if n.regex.String() == regex.String() {
			return n, true
		}
	}

	return nil, false
}

type node struct {
	regex    *regexp.Regexp
	handlers HandlersChain
}

func (n *node) addHandlers(handlers HandlersChain) {
	newHandlers := append(n.handlers, handlers...)
	if int8(len(newHandlers)) > abortIndex {
		panic("so more handlers")
	}
	n.handlers = newHandlers
}
