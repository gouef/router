package router

import (
	"github.com/Gouef/router/http"
	"net/url"
)

type RouteList struct {
	list []struct {
		Router Router
		Rank   int
	}

	ranks  [][]Router
	domain *string
	path   *string
}

func NewRouteList() *RouteList {
	return &RouteList{}
}

func (r RouteList) constructUrl(params map[string]interface{}, refUrl http.Url) *string {

}

func (r RouteList) match(httpRequest http.IRequest) *url.Values {

}
