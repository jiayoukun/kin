package kin

import (
	"encoding/json"
	"fmt"
	"kin/binding"
	"net/http"
	"sync"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Path string
	Method string
	Params map[string]string
	StatusCode	int
	index int
	handlers []HandlerFunc
	mu sync.RWMutex
	Keys map[string]interface{}
}

func newContext(w http.ResponseWriter,req *http.Request) *Context {
	return &Context{
		Path:	req.URL.Path,
		Method:	req.Method,
		Req:	req,
		Writer:  w,
		index: -1,
	}
}

func (c*Context)Next()  {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++{
		c.handlers[c.index](c)
	}
}

func (c*Context)Param(key string) string {
	value, _ :=	c.Params[key]
	return value
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c*Context)Status(code int)  {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context)SetHeader(key string, value string)  {
	c.Writer.Header().Set(key,value)
}
func (c*Context)String(code int, format string,value ...interface{})  {
	c.SetHeader("Content-Type","text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context)JSON(code int, obj interface{})  {
	c.SetHeader("Content-Type","application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj);err != nil{
		http.Error(c.Writer,err.Error(),500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.GetHeader("Content-Type"))
}

// Get the value of Request-Header
func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

// Bind data to struct
func (c *Context) Bind(obj interface{}) error {
	b := binding.Defaultt(c.ContentType())
	return c.BindData(obj, b)
}

// Call the bind() of the binding instance
func (c *Context) BindData(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Req, obj)
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}
