package tests

import (
	"github.com/gin-gonic/gin"
	router2 "github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterErrorHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := router2.NewRouter()

	router.SetErrorHandler(http.StatusNotFound, func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom 404",
		})
	})

	router.SetErrorHandler(http.StatusInternalServerError, func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Custom 500",
		})
	})

	router.AddRouteGet("ok", "/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is OK"})
	})
	router.AddRouteGet("notfound", "/notfound", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Custom 404",
		})
	})
	router.AddRouteGet("servererror", "/servererror", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Custom 500",
		})
	})

	// Test: `/ok` (status 200)
	t.Run("Test OK Status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		w := httptest.NewRecorder()

		router.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"This is OK"}`, w.Body.String())
	})

	// Test: `/notfound` (status 404)
	t.Run("Test Custom 404 Handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		w := httptest.NewRecorder()

		router.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Custom 404"}`, w.Body.String())
	})

	// Test: `/servererror` (status 500)
	t.Run("Test Custom 500 Handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/servererror", nil)
		w := httptest.NewRecorder()

		router.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Custom 500"}`, w.Body.String())
	})

	// Test: unknown route (status 404)
	t.Run("Test Unknown Route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
		w := httptest.NewRecorder()

		router.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Custom 404"}`, w.Body.String())
	})
}
