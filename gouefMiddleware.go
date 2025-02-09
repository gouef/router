package router

import "github.com/gin-gonic/gin"

func GouefMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Powered-By", "Gouef")
		c.Next()
	}
}
