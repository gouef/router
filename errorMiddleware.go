package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Pokračuj v zpracování požadavku

		// Získání statusu odpovědi
		statusCode := c.Writer.Status()

		// Pokud je status chybový (4xx nebo 5xx), přepíšeme odpověď
		if statusCode >= 400 {
			switch statusCode {
			case http.StatusNotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error":       "Resource not found",
					"description": "The requested endpoint does not exist on the server",
				})
			case http.StatusInternalServerError:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       "Internal server error",
					"description": "An unexpected error occurred on the server",
				})
			default:
				c.JSON(statusCode, gin.H{
					"error":       "An error occurred",
					"description": "Please check the documentation for this API",
				})
			}
		}
	}
}