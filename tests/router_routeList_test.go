package tests

import (
	"github.com/gin-gonic/gin"
	router2 "github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouteList(t *testing.T) {
	type ProductDetailParams struct {
		Locale string `uri:"locale" binding:"required"`
		ID     int    `uri:"id" binding:"required"`
	}
	productDetailHandler := func(c *gin.Context, params *ProductDetailParams) {
		c.JSON(http.StatusOK, gin.H{
			"locale": params.Locale,
			"id":     params.ID,
		})

		assert.Equal(t, "cs", params.Locale)
		assert.Equal(t, 42, params.ID)
	}

	lr := router2.NewRouteList()
	v1 := router2.CreateRouteList("/v1")
	lr.AddChild(v1)

	lr.Add("products:detail", "/:locale/products/:id", productDetailHandler, router2.Get)
	v1.Add("products:detail", "/:locale/products/:id", productDetailHandler, router2.Get)

	router := router2.NewRouter()
	router.AddRouteList(lr)
	router2.CreateRoute(router, "test", "/test/:id", func(c *gin.Context, p *struct {
		ID int `uri:"id" binding:"required"`
	}) {
		c.JSON(http.StatusOK, gin.H{
			"id": p.ID,
		})
		assert.Equal(t, 42, p.ID)
	}, router2.Get)

	// Testování požadavku
	req := httptest.NewRequest(http.MethodGet, "/cs/products/42", nil)
	w := httptest.NewRecorder()
	router.GetNativeRouter().ServeHTTP(w, req)

	// Ověření výsledku
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"locale":"cs","id":42}`, w.Body.String())

	// Testování požadavku
	req2 := httptest.NewRequest(http.MethodGet, "/v1/cs/products/42", nil)
	w2 := httptest.NewRecorder()
	router.GetNativeRouter().ServeHTTP(w2, req2)

	// Ověření výsledku
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.JSONEq(t, `{"locale":"cs","id":42}`, w2.Body.String())

}
