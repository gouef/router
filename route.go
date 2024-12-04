package router

type Route struct {
	path   string
	method Method
	handle func()
}

func NewRoute(path string, handle func()) *Route {
	return &Route{
		path:   path,
		handle: handle,
	}
}

func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) GetHandle() func() {
	return r.handle
}
