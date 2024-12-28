package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetters(t *testing.T) {
	children := map[string]*router.Route{}
	handler := func(c *gin.Context) {}
	pattern := "/test"
	route := router.NewRoute("test", pattern, handler, router.Get, children)
	assert.Equal(t, router.Get, route.GetMethod())
	assert.Equal(t, pattern, route.GetPattern())
}
