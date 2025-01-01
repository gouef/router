package tests

import (
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMethodFromString(t *testing.T) {
	m, _ := router.MethodFromString("get")
	assert.Equal(t, router.Get, m)
}

func TestMethodString(t *testing.T) {
	const TEST router.Method = 30
	assert.Equal(t, "GET", router.Get.String())
	assert.Equal(t, "PUT", router.Put.String())
	assert.Equal(t, "PATCH", router.Patch.String())
	assert.Equal(t, "POST", router.Post.String())
	assert.Equal(t, "DELETE", router.Delete.String())
	assert.Equal(t, "OPTIONS", router.Options.String())
	assert.Equal(t, "CONNECT", router.Connect.String())
	assert.Equal(t, "HEAD", router.Head.String())
	assert.Equal(t, "TRACE", router.Trace.String())
	assert.Equal(t, "Unknown", TEST.String())
}
