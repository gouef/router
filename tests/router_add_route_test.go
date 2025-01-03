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

func TestAddRouteError(t *testing.T) {
	t.Run("checkParamType", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()
		handler := func(c *gin.Context, param string) {}
		err := r.AddRoute("error", "/error", handler, router.Get)
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.NotNil(t, err)
		assert.Equal(t, `Handler parameter must be a pointer to a struct, got string`, err.Error())
	})

	t.Run("first param must be a function", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()
		//handler := func(c *gin.Context, param string) {}
		err := r.AddRoute("error", "/error", "string", router.Get)
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.NotNil(t, err)
		assert.Equal(t, `handler must be a function`, err.Error())
	})

	t.Run("function has more parameters", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()
		handler := func(c *gin.Context, param string, param2 string) {}
		err := r.AddRoute("error", "/error", handler, router.Get)
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.NotNil(t, err)
		assert.Equal(t, `handler must have one or two parameters`, err.Error())
	})

	t.Run("Error createNativeRoute", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()

		handler := func(c *gin.Context, param string, param2 string) {}

		route2 := router.NewRoute("error2", "/error2", handler, router.Get)
		rl := router.CreateRouteList("")

		rl.AddRoute(route2)

		err := r.AddRouteList(rl)

		req := httptest.NewRequest(http.MethodGet, "/v1/error2", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, `handler must have one or two parameters`, err.Error())
		}
	})

	t.Run("Error createNativeRoute Children", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()

		type TestParam struct {
			ID int `uri:"id"`
		}

		mHandler := func(c *gin.Context, param *TestParam) {}
		handler := func(c *gin.Context, param string, param2 string) {}

		route := router.NewRoute("error", "/error", mHandler, router.Get)
		route2 := router.NewRoute("error2", "/error2", handler, router.Get)
		rl := router.CreateRouteList("")
		rlV1 := router.CreateRouteList("/v1")
		rlV1E := router.CreateRouteList("/errors")

		rlV1.AddRoute(route)
		rlV1E.AddRoute(route2)
		rlV1.AddChild(rlV1E)
		rl.AddChild(rlV1)

		err := r.AddRouteList(rl)

		req := httptest.NewRequest(http.MethodGet, "/v1/error2", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.NotNil(t, err)
		if err != nil {
			assert.Equal(t, `handler must have one or two parameters`, err.Error())
		}
	})

	t.Run("bind parameters error", func(t *testing.T) {

		r := router.NewRouter()
		r.EnableTest()
		type TestParam struct {
			ID int `uri:"id" binding:"required"`
		}
		handler := func(c *gin.Context, param *TestParam) {
			c.JSON(http.StatusOK, gin.H{"id": param.ID})
		}
		err := r.AddRoute("error", "/error", handler, router.Get)
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		r.GetNativeRouter().ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, "{\"error\": \"Key: 'TestParam.ID' Error:Field validation for 'ID' failed on the 'required' tag\"}", w.Body.String())
		assert.Nil(t, err)
	})
}
