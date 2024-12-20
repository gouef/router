package tests

import (
	"github.com/gin-gonic/gin"
	router2 "github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	router := router2.NewRouter()

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
	router.AddRoute("products:detail", "/:locale/products/:id", productDetailHandler, router2.Get)
	router.AddRoute(
		"product:detail",
		"/product/:id",
		func(c *gin.Context, p *struct {
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
	req2 := httptest.NewRequest(http.MethodGet, "/product/42", nil)
	w2 := httptest.NewRecorder()
	router.GetNativeRouter().ServeHTTP(w2, req2)

	// Ověření výsledku
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.JSONEq(t, `{"id":42}`, w2.Body.String())
}

func TestRouter_AddRouteWithoutParams(t *testing.T) {
	router := router2.NewRouter()

	productDetailHandler := func(c *gin.Context) {
		l := c.Param("locale")
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"locale": l,
			"id":     id,
		})

		assert.Equal(t, "cs", l)
		assert.Equal(t, "42", id)
	}

	// Přidání routy
	router.AddRoute("products:detail", "/:locale/products/:id", productDetailHandler, router2.Get)

	// Testování požadavku
	req := httptest.NewRequest(http.MethodGet, "/cs/products/42", nil)
	w := httptest.NewRecorder()
	router.GetNativeRouter().ServeHTTP(w, req)

	// Ověření výsledku
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"locale":"cs","id":"42"}`, w.Body.String())
}

func TestRouterRoutesKeys(t *testing.T) {
	// Vytvoření nového routeru
	r := router2.NewRouter()

	// Přidání tras
	r.AddRoute("home", "/home", func(c *gin.Context) {}, router2.Get)
	r.AddRoute("about", "/about", func(c *gin.Context) {}, router2.Get)
	r.AddRoute("contact", "/contact", func(c *gin.Context) {}, router2.Get)
	r.AddRoute("product:detail", "/product/:id", func(c *gin.Context) {}, router2.Get)

	// Očekávané klíče
	expectedKeys := []string{"home", "about", "contact", "product:detail"}

	// Kontrola klíčů
	for _, key := range expectedKeys {
		if _, exists := r.GetRoutes()[key]; !exists {
			t.Errorf("Expected key '%s' to exist in Router.routes, but it does not.", key)
		}
	}

	// Kontrola počtu klíčů
	if len(r.GetRoutes()) != len(expectedKeys) {
		t.Errorf("Expected %d routes, but got %d", len(expectedKeys), len(r.GetRoutes()))
	}
}

func TestRoute_AddChild(t *testing.T) {
	// Vytvoření root route
	root := router2.NewRoute("root", "/", nil, router2.Get, make(map[string]*router2.Route))

	// Přidání dětí
	root.AddChild(
		"child1",
		"/child1",
		nil,
		router2.Get,
	).AddChild(
		"child2",
		"/child2",
		nil,
		router2.Post,
	)

	// Ověření, že root má 2 děti
	assert.Len(t, root.GetChildren(), 2, "Root should have 2 children")

	// Ověření existence klíčů pro děti
	assert.Contains(t, root.GetChildren(), "/child1", "Expected '/child1' in children")
	assert.Contains(t, root.GetChildren(), "/child2", "Expected '/child2' in children")

	// Ověření správné instance dítěte
	assert.Equal(t, "root:child1", root.GetChildren()["/child1"].GetName(), "Expected the same child1 instance")
	assert.Equal(t, "root:child2", root.GetChildren()["/child2"].GetName(), "Expected the same child2 instance")
}
