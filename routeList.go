package router

type RouteList struct {
	pattern  string
	routes   []*Route
	children []*RouteList
}

func NewRouteList() *RouteList {
	return CreateRouteList("")
}

func CreateRouteList(pattern string) *RouteList {
	return &RouteList{pattern: pattern}
}

func (l *RouteList) AddChild(child *RouteList) *RouteList {
	l.children = append(l.children, child)

	return l
}

func (l *RouteList) Add(name string, pattern string, handler interface{}, method Method) *RouteList {
	return l.AddWithChildren(name, pattern, handler, method, nil)
}

func (l *RouteList) AddWithChildren(name string, pattern string, handler interface{}, method Method, children map[string]*Route) *RouteList {
	route := NewRoute(name, pattern, handler, method, children)
	l.AddRoute(route)

	return l
}

func (l *RouteList) AddRoute(r *Route) *RouteList {
	l.routes = append(l.routes, r)
	return l
}
