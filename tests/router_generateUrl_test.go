package tests

import (
	"github.com/gin-gonic/gin"
	router "github.com/gouef/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateUrlByName(t *testing.T) {
	r := router.NewRouter()
	r.AddRouteObject(router.NewRoute("index", "/", func(c *gin.Context) {}, router.Get, nil))
	r.AddRouteObject(router.NewRoute("home", "/home", func(c *gin.Context) {}, router.Get, nil))
	r.AddRouteObject(router.NewRoute("user", "/user/:id", func(c *gin.Context) {}, router.Get, nil))

	tests := []struct {
		name       string
		params     map[string]interface{}
		expected   string
		shouldFail bool
	}{
		// Test 1: Generování URL pro existující route "index"
		{
			name:     "index",
			params:   nil,
			expected: "/",
		},
		// Test 2: Generování URL pro existující route "home"
		{
			name:     "home",
			params:   nil,
			expected: "/home",
		},
		// Test 3: Generování URL pro existující route "user" s parametrem
		{
			name:     "user",
			params:   map[string]interface{}{"id": 42},
			expected: "/user/42",
		},
		// Test 4: Chyba při neexistující route "unknown"
		{
			name:       "unknown",
			params:     nil,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.GenerateUrlByName(tt.name, tt.params)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGenerateUrlByPattern(t *testing.T) {
	r := router.NewRouter()

	tests := []struct {
		pattern    string
		params     map[string]interface{}
		expected   string
		shouldFail bool
	}{
		// Test 1: Generování URL s parametrem :id
		{
			pattern:  "/user/:id",
			params:   map[string]interface{}{"id": 42},
			expected: "/user/42",
		},
		// Test 2: Chyba při chybějícím parametru :id
		{
			pattern:    "/user/:id",
			params:     nil,
			shouldFail: true,
		},
		// Test 3: Generování URL bez parametrů
		{
			pattern:  "/home",
			params:   nil,
			expected: "/home",
		},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			result, err := r.GenerateUrlByPattern(tt.pattern, tt.params)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
