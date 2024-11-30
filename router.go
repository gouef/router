package router

import (
	"github.com/Gouef/router/http"
	"net/url"
)

type Router interface {
	match(httpRequest http.IRequest) *url.Values
	constructUrl(params []any, refUrl http.Url) *string
}

const ONE_WAY = true
