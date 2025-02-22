package router

import "github.com/gin-gonic/gin"

type HandlerFunc func(response *Response, request *Request)
type HandlerParamFunc[T any] func(response *Response, request *Request, p *T)

type HandlerContext func(c *gin.Context)

type HandlerContextParam[T any] func(c *gin.Context, p *T)

type Route struct {
	Name    string
	Pattern string
	Handler interface{}
	Method  Method
}

// NewRoute create Route
func NewRoute(name string, pattern string, handler interface{}, method Method) *Route {
	return &Route{
		Name:    name,
		Pattern: pattern,
		Handler: handler,
		Method:  method,
	}
}

// GetName get route Name
func (r *Route) GetName() string {
	return r.Name
}

// GetMethod get route Method
func (r *Route) GetMethod() Method {
	return r.Method
}

// GetPattern get route Pattern
func (r *Route) GetPattern() string {
	return r.Pattern
}

// GetHandler get route Handler
func (r *Route) GetHandler() interface{} {
	return r.Handler
}
