package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/router/http"
	"net/url"
)

type Router interface {
	match(httpRequest http.IRequest) *url.Values
	constructUrl(params map[string]interface{}, refUrl http.Url) *string
}

const ONE_WAY = true

func run() {
	router := gin.Default()

	router.Handle()
}
