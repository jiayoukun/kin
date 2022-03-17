package kin

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Handler interface {
	ServeHTTP(w http.ResponseWriter,req *http.Request)
}

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

func Default() *Engine {
	engine	:= New()
	engine.Use(Logger(), Recovery())
	return engine
}

func New() *Engine {
	engine :=&Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (engine *Engine)Handlers() http.HandlerFunc {
	return engine.Handlers()
}

func (engine *Engine)Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine.Handlers())
}



func (engine *Engine)ServerHTTP(w http.ResponseWriter,req *http.Request)  {
	var middlewares	[]HandlerFunc
	for _, group := range engine.groups{
		if strings.HasPrefix(req.URL.Path, group.prefix){
			middlewares = append(middlewares,group.middlewares...)
		}
	}
	c := newContext(w,req)
	c.handlers = middlewares
	engine.router.handle(c)
}
