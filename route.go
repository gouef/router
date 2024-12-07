package router

type Route struct {
	pattern  string
	handler  interface{}
	method   Method
	children map[string]*Route
}

func NewRoute(pattern string, handler interface{}, method Method, children map[string]*Route) *Route {
	return &Route{
		pattern:  pattern,
		handler:  handler,
		method:   method,
		children: children,
	}
}

func (r *Route) AddChild(pattern string, handler interface{}, method Method) *Route {
	child := NewRoute(pattern, handler, method, nil)
	r.children[pattern] = child

	return r
}

func (r *Route) GetChildren() map[string]*Route {
	return r.children
}

func (r *Route) GetPattern() string {
	return r.pattern
}

func (r *Route) GetHandler() interface{} {
	return r.handler
}
