package kin

import "log"

type RouterGroup struct {
	prefix		string
	middlewares	[]HandlerFunc
	parent 		*RouterGroup
	engine 		*Engine
}

func (g*RouterGroup)Group(prefix string)*RouterGroup  {
	engine :=	g.engine
	newGroup := &RouterGroup{
		prefix:	g.prefix + prefix,
		parent: g.parent,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
