package router

import "github.com/gin-gonic/gin"

func (r *Router) IsDebug() bool {
	return IsDebug()
}

func (r *Router) EnableDebug() {
	EnableDebug()
}
func (r *Router) IsTest() bool {
	return IsTest()
}

func (r *Router) EnableTest() {
	EnableTest()
}
func (r *Router) IsRelease() bool {
	return IsRelease()
}

func (r *Router) EnableRelease() {
	EnableRelease()
}

func (r *Router) SetMode(mode string) *Router {
	gin.SetMode(mode)
	return r
}

func SetMode(mode string) {
	gin.SetMode(mode)
}

func GetMode() string {
	return gin.Mode()
}

func (r *Router) GetMode() string {
	return GetMode()
}

func IsDebug() bool {
	return gin.IsDebugging()
}

func EnableDebug() {
	SetMode(DebugMode)
}

func IsTest() bool {
	return gin.Mode() == TestMode
}

func EnableTest() {
	SetMode(TestMode)
}

func IsRelease() bool {
	return gin.Mode() == ReleaseMode
}

func EnableRelease() {
	SetMode(ReleaseMode)
}
