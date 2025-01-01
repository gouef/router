package router

import "net/http"

type Method int

const (
	Get Method = iota
	Post
	Head
	Put
	Delete
	Patch
	Options
	Connect
	Trace
)

var stringToMethod = map[string]Method{
	http.MethodGet:     Get,
	http.MethodPost:    Post,
	http.MethodHead:    Head,
	http.MethodPut:     Put,
	http.MethodDelete:  Delete,
	http.MethodPatch:   Patch,
	http.MethodOptions: Options,
	http.MethodConnect: Connect,
	http.MethodTrace:   Trace,
}

func (m Method) String() string {
	switch m {
	case Get:
		return http.MethodGet
	case Post:
		return http.MethodPost
	case Head:
		return http.MethodHead
	case Put:
		return http.MethodPut
	case Delete:
		return http.MethodDelete
	case Patch:
		return http.MethodPatch
	case Options:
		return http.MethodOptions
	case Connect:
		return http.MethodConnect
	case Trace:
		return http.MethodTrace
	default:
		return "Unknown"
	}
}

func MethodFromString(method string) (Method, bool) {
	m, exists := stringToMethod[method]
	return m, exists
}
