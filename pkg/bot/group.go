package bot

import "regexp"

type Router interface {
	Use(...HandlerFunc) Router

	Command(string, ...HandlerFunc) Router
	Message(string, ...HandlerFunc) Router
	Callback(string, ...HandlerFunc) Router

	CommandRegex(*regexp.Regexp, ...HandlerFunc) Router
	MessageRegex(*regexp.Regexp, ...HandlerFunc) Router
	CallbackRegex(*regexp.Regexp, ...HandlerFunc) Router

	Group(string, HandlerGroup) Router
	GroupState(string, HandlerGroup) Router
	Handle(string, *regexp.Regexp, ...HandlerFunc) Router
}

type RouterGroup struct {
	engine *Engine

	state     string
	startPath string
}

type HandlerGroup func(group Router)

func (group *RouterGroup) Use(middlewares ...HandlerFunc) Router {
	group.handle("*", regexForAllStrings, middlewares)
	return group
}

func (group *RouterGroup) Command(template string, handlers ...HandlerFunc) Router {
	group.handle("command", group.engine.mustCompileRegexByTemplate(template), handlers)
	return group
}

func (group *RouterGroup) Message(template string, handlers ...HandlerFunc) Router {
	group.handle("message", group.engine.mustCompileRegexByTemplate(template), handlers)
	return group
}

func (group *RouterGroup) Callback(template string, handlers ...HandlerFunc) Router {
	group.handle("callback", group.engine.mustCompileRegexByTemplate(template), handlers)
	return group
}

func (group *RouterGroup) Group(startPath string, handlerGroup HandlerGroup) Router {
	routerGroup := &RouterGroup{engine: group.engine, startPath: startPath}
	handlerGroup(routerGroup)

	return routerGroup
}

func (group *RouterGroup) GroupState(state string, handlerGroup HandlerGroup) Router {
	routerGroup := &RouterGroup{engine: group.engine, state: state}
	handlerGroup(routerGroup)

	return routerGroup
}

func (group *RouterGroup) CommandRegex(regex *regexp.Regexp, handlers ...HandlerFunc) Router {
	group.handle("command", regex, handlers)
	return group
}

func (group *RouterGroup) MessageRegex(regex *regexp.Regexp, handlers ...HandlerFunc) Router {
	group.handle("message", regex, handlers)
	return group
}

func (group *RouterGroup) CallbackRegex(regex *regexp.Regexp, handlers ...HandlerFunc) Router {
	group.handle("callback", regex, handlers)
	return group
}

func (group *RouterGroup) Handle(method string, regex *regexp.Regexp, handlers ...HandlerFunc) Router {
	group.handle(method, regex, handlers)
	return group
}

func (group *RouterGroup) handle(method string, regex *regexp.Regexp, handlers HandlersChain) {
	tree, ok := group.engine.trees.get(method)
	if !ok {
		tree = &methodTree{method: method, nodes: make([]*node, 0)}
		group.engine.trees = append(group.engine.trees, tree)
	}

	root, ok := tree.get(regex, group.state, group.startPath)
	if !ok {
		root = new(node)
		root.regex = regex
		root.state = group.state
		root.startPath = group.startPath

		tree.nodes = append(tree.nodes, root)
	}
	root.addHandlers(handlers)
}
