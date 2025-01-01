package tests

import (
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMode(t *testing.T) {
	r := router.NewRouter()

	assert.Equal(t, false, r.IsMode("release"))
	assert.Equal(t, true, r.IsMode("debug"))
	assert.Equal(t, false, r.IsMode("something"))
}

func TestEnableMode(t *testing.T) {
	r := router.NewRouter()

	assert.Equal(t, true, r.EnableMode("release"))
	assert.True(t, r.IsRelease())
	assert.True(t, r.EnableRelease())
	assert.True(t, r.IsRelease())
	assert.Equal(t, true, r.EnableMode("test"))
	assert.True(t, r.IsTest())
	assert.True(t, r.EnableTest())
	assert.True(t, r.IsTest())
	assert.Equal(t, true, r.EnableMode("debug"))
	assert.True(t, r.IsDebug())
	assert.True(t, r.EnableDebug())
	assert.True(t, r.IsDebug())
	assert.Equal(t, false, r.EnableMode("something"))
}

func TestSetMode(t *testing.T) {
	r := router.NewRouter()

	assert.Equal(t, true, r.SetMode("release").IsRelease())
	assert.Equal(t, "release", r.GetMode())
}

func TestGinMode(t *testing.T) {
	router.EnableRelease()
	assert.Equal(t, "release", router.GetMode())
	assert.True(t, router.IsRelease())
	assert.False(t, router.IsDebug())
	router.EnableTest()
	assert.Equal(t, "test", router.GetMode())
	assert.True(t, router.IsTest())
	assert.False(t, router.IsRelease())
	router.EnableDebug()
	assert.Equal(t, "debug", router.GetMode())
	assert.True(t, router.IsDebug())
	assert.False(t, router.IsTest())
}
