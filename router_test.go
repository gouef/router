package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	router := NewRouter()

	// DTO a handler
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

	// Přidání routy
	router.AddRoute("/:locale/products/:id", productDetailHandler, Get)

	// Testování požadavku
	req := httptest.NewRequest(http.MethodGet, "/cs/products/42", nil)
	w := httptest.NewRecorder()
	router.GetNativeRouter().ServeHTTP(w, req)

	// Ověření výsledku
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"locale":"cs","id":42}`, w.Body.String())
}
