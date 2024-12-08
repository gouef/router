package router

type Route struct {
	name     string
	pattern  string
	handler  interface{}
	method   Method
	children map[string]*Route
}

func NewRoute(name string, pattern string, handler interface{}, method Method, children map[string]*Route) *Route {
	return &Route{
		name:     name,
		pattern:  pattern,
		handler:  handler,
		method:   method,
		children: children,
	}
}

func (r *Route) AddChild(name string, pattern string, handler interface{}, method Method) *Route {
	child := NewRoute(r.name+":"+name, pattern, handler, method, nil)
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
