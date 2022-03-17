package kin

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func NewRouter()*router  {
	return &router{
		roots:	make(map[string]*node),
		handlers: 	make(map[string]HandlerFunc),
	}
}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern,"/")
	parts := make([]string,0)
	for _, v := range vs{
		if v != ""{
			parts = append(parts,v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router)addRoute(method string,patterm string, handler HandlerFunc)  {
	parts := parsePattern(patterm)
	key :=	method + patterm
	_,ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	_, ok = r.roots[method].judge(parts,0)
	if !ok {
		panic(patterm + ": The route conflicts with the previous route")
	}
	r.roots[method].insert(patterm,parts,0)
	r.handlers[key] = handler
}

func (r *router)getRoute(method string,path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root,ok := r.roots[method]
	if !ok {
		return nil,nil
	}
	n := root.search(searchParts,0)
	if n != nil{
		parts := parsePattern(n.pattern)
		for index, part := range parts{
			if	part[0] == ':' {
				params[part[1:]] =searchParts[index]
			}
			if part[0] == '*'&& len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:],"/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
func (r *router)handle(c *Context)  {
	n,params :=	r.getRoute(c.Method,c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers,r.handlers[key])
	}else{
		c.handlers = append(c.handlers, func(context *Context) {
			c.String(http.StatusNotFound,"404 NOT FOUND: %s\n",c.Path)
		})
	}
	c.Next()
}
