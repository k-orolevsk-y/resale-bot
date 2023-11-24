package bot

import (
	"regexp"
	"strings"
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

func (trees methodTrees) compileHandlersChain(method, query, state string) HandlersChain {
	var handlers HandlersChain
	onlyMiddlewares := true

	for _, tree := range trees {
		if tree.method == method || tree.method == "*" {
			for _, n := range tree.nodes {
				if n.startPath != "" && !strings.HasPrefix(query, n.startPath) {
					continue
				}

				if n.state != "" && n.state != state {
					continue
				}

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

func (tree *methodTree) get(regex *regexp.Regexp, state, startPath string) (*node, bool) {
	for _, n := range tree.nodes {
		if n.regex.String() == regex.String() && n.state == state && n.startPath == startPath {
			return n, true
		}
	}

	return nil, false
}

type node struct {
	regex     *regexp.Regexp
	state     string
	startPath string

	handlers HandlersChain
}

func (n *node) addHandlers(handlers HandlersChain) {
	newHandlers := append(n.handlers, handlers...)
	if int8(len(newHandlers)) > abortIndex {
		panic("so more handlers")
	}
	n.handlers = newHandlers
}
