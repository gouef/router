package router

type RouteList struct {
	pattern  string
	routes   []*Route
	children []*RouteList
}

// NewRouteList create RouteList with empty pattern
func NewRouteList() *RouteList {
	return CreateRouteList("")
}

// CreateRouteList create RouteList
func CreateRouteList(pattern string) *RouteList {
	return &RouteList{pattern: pattern}
}

// AddChild add child of RouteList
func (l *RouteList) AddChild(child *RouteList) *RouteList {
	l.children = append(l.children, child)

	return l
}

// Add create route for RouteList
func (l *RouteList) Add(name string, pattern string, handler interface{}, method Method) *RouteList {
	route := NewRoute(name, pattern, handler, method)
	l.AddRoute(route)

	return l
}

// AddRoute add Route to RouteList
func (l *RouteList) AddRoute(r *Route) *RouteList {
	l.routes = append(l.routes, r)
	return l
}
