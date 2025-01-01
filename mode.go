package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/mode"
)

func (r *Router) IsMode(mode string) bool {
	m, err := r.Mode.IsMode(mode)

	if err != nil {
		return false
	}

	return m
}

func (r *Router) EnableMode(mode string) bool {
	m, err := r.Mode.EnableMode(mode)

	if err != nil {
		return false
	}

	switch mode {
	case DebugMode:
		EnableDebug()
	case ReleaseMode:
		EnableRelease()
	default:
		EnableTest()
	}
	return m
}

func (r *Router) IsDebug() bool {
	return r.IsMode(mode.DebugMode)
}

func (r *Router) EnableDebug() bool {
	return r.EnableMode(mode.DebugMode)
}
func (r *Router) IsTest() bool {
	return r.IsMode(mode.TestMode)
}

func (r *Router) EnableTest() bool {
	return r.EnableMode(mode.TestMode)
}
func (r *Router) IsRelease() bool {
	return r.IsMode(mode.ReleaseMode)
}

// EnableRelease set Mode release for project and gin
func (r *Router) EnableRelease() bool {
	return r.EnableMode(mode.ReleaseMode)
}

// SetMode set Mode for project and gin
func (r *Router) SetMode(mode string) *Router {
	r.EnableMode(mode)
	return r
}

// SetMode set Mode for gin
func SetMode(mode string) {
	gin.SetMode(mode)
}

// GetMode return gin Mode
func GetMode() string {
	return gin.Mode()
}

// GetMode return project Mode
func (r *Router) GetMode() string {
	return r.Mode.GetMode()
}

// IsDebug if gin Mode is debug
func IsDebug() bool {
	return gin.IsDebugging()
}

// EnableDebug set Mode debug for gin
func EnableDebug() {
	SetMode(DebugMode)
}

// IsTest if gin Mode is test
func IsTest() bool {
	return gin.Mode() == TestMode
}

// EnableTest set Mode test for gin
func EnableTest() {
	SetMode(TestMode)
}

// IsRelease if gin Mode is release
func IsRelease() bool {
	return gin.Mode() == ReleaseMode
}

// EnableRelease set Mode release for gin
func EnableRelease() {
	SetMode(ReleaseMode)
}
