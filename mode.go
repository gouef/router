package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/mode"
)

func (r *Router) isMode(mode string) bool {
	m, err := r.mode.IsMode(mode)

	if err != nil {
		panic("non exists mode")
	}

	return m
}

func (r *Router) EnableMode(mode string) bool {
	m, err := r.mode.EnableMode(mode)

	if err != nil {
		panic("non exists mode")
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
	return r.isMode(mode.DebugMode)
}

func (r *Router) EnableDebug() bool {
	return r.EnableMode(mode.DebugMode)
}
func (r *Router) IsTest() bool {
	return r.isMode(mode.TestMode)
}

func (r *Router) EnableTest() bool {
	return r.EnableMode(mode.TestMode)
}
func (r *Router) IsRelease() bool {
	return r.isMode(mode.ReleaseMode)
}

// EnableRelease set mode release for project and gin
func (r *Router) EnableRelease() bool {
	return r.EnableMode(mode.ReleaseMode)
}

// SetMode set mode for project and gin
func (r *Router) SetMode(mode string) *Router {
	r.EnableMode(mode)
	return r
}

// SetMode set mode for gin
func SetMode(mode string) {
	gin.SetMode(mode)
}

// GetMode return gin mode
func GetMode() string {
	return gin.Mode()
}

// GetMode return project mode
func (r *Router) GetMode() string {
	return r.mode.GetMode()
}

// IsDebug if gin mode is debug
func IsDebug() bool {
	return gin.IsDebugging()
}

// EnableDebug set mode debug for gin
func EnableDebug() {
	SetMode(DebugMode)
}

// IsTest if gin mode is test
func IsTest() bool {
	return gin.Mode() == TestMode
}

// EnableTest set mode test for gin
func EnableTest() {
	SetMode(TestMode)
}

// IsRelease if gin mode is release
func IsRelease() bool {
	return gin.Mode() == ReleaseMode
}

// EnableRelease set mode release for gin
func EnableRelease() {
	SetMode(ReleaseMode)
}
