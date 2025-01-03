package router

type Route struct {
	name    string
	pattern string
	handler interface{}
	method  Method
}

// NewRoute create Route
func NewRoute(name string, pattern string, handler interface{}, method Method) *Route {
	return &Route{
		name:    name,
		pattern: pattern,
		handler: handler,
		method:  method,
	}
}

// GetName get route name
func (r *Route) GetName() string {
	return r.name
}

// GetMethod get route method
func (r *Route) GetMethod() Method {
	return r.method
}

// GetPattern get route pattern
func (r *Route) GetPattern() string {
	return r.pattern
}

// GetHandler get route handler
func (r *Route) GetHandler() interface{} {
	return r.handler
}
