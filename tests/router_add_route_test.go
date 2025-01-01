package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := router.NewRouter()

	r.AddRouteGet("get", "/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
	})

	r.AddRouteMethod("get2", "/get2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get2",
		})
	}, router.Get)

	r.AddMultiMethodsRoute("users", "/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "users",
		})
	}, []router.Method{router.Get, router.Post})

	r.AddRoutePost("post", "/post", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "post",
		})
	})

	r.AddRoutePatch("patch", "/patch", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "patch",
		})
	})

	r.AddRoutePut("put", "/put", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "put",
		})
	})

	r.AddRouteDelete("delete", "/delete", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "delete",
		})
	})

	r.AddRouteHead("head", "/head", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "head",
		})
	})

	r.AddRouteOptions("options", "/options", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "options",
		})
	})

	r.AddRouteTrace("trace", "/trace", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "trace",
		})
	})

	r.AddRouteConnect("connect", "/connect", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "connect",
		})
	})

	t.Run("AddMultiRouteMethod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"users"}`, w.Body.String())

		req = httptest.NewRequest(http.MethodPost, "/users", nil)
		w = httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"users"}`, w.Body.String())
	})

	t.Run("AddRouteMethod", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get2", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"get2"}`, w.Body.String())
	})

	t.Run("AddRouteGet", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"get"}`, w.Body.String())
	})

	t.Run("AddRoutePost", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/post", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"post"}`, w.Body.String())
	})

	t.Run("AddRoutePatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/patch", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"patch"}`, w.Body.String())
	})

	t.Run("AddRouteDelete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"delete"}`, w.Body.String())
	})

	t.Run("AddRoutePut", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/put", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"put"}`, w.Body.String())
	})

	t.Run("AddRouteHead", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/head", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"head"}`, w.Body.String())
	})

	t.Run("AddRouteOptions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/options", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"options"}`, w.Body.String())
	})

	t.Run("AddRouteConnect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodConnect, "/connect", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"connect"}`, w.Body.String())
	})

	t.Run("AddRouteTrace", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodTrace, "/trace", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"trace"}`, w.Body.String())
	})

}
