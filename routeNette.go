package router

import (
	"github.com/gouef/router/http"
	"net/url"
)

type RouteNette struct {
	mask     string
	sequence []interface{}

	/** regular expression pattern */
	re string

	/** @var string[]  parameter aliases in regular expression */
	aliases []string

	/** @var array of [value & fixity, filterIn, filterOut] */
	metadata []interface{}
	xlat     []interface{}

	/** Host, Path, Relative */
	routeType int

	/** http | https */
	scheme string
}

func NewRouteNette() *RouteNette {
	return &RouteNette{}
}

func (r RouteNette) match(httpRequest http.IRequest) *url.Values {

}

func (r RouteNette) constructUrl(params map[string]interface{}, refUrl http.Url) *string {

}
