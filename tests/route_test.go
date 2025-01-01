package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetters(t *testing.T) {
	children := map[string]*router.Route{}
	called := false
	handler := func(c *gin.Context) {
		called = true
	}
	pattern := "/test"
	route := router.NewRoute("test", pattern, handler, router.Get, children)

	handlerFromRoute, ok := route.GetHandler().(func(*gin.Context))
	assert.True(t, ok, "Handler is not of type func(*gin.Context)")

	c := &gin.Context{}
	handlerFromRoute(c)

	assert.Equal(t, router.Get, route.GetMethod())
	assert.Equal(t, pattern, route.GetPattern())
	assert.True(t, called)
}
