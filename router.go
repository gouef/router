package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ErrorHandlerFunc func(c *gin.Context)

type Router struct {
	router         *gin.Engine
	routes         map[string]*Route
	errorHandlers  map[int]ErrorHandlerFunc
	defaultHandler ErrorHandlerFunc
}

// NewRouter create new Router
func NewRouter() *Router {
	router := gin.New()
	return &Router{
		router:        router,
		routes:        make(map[string]*Route),
		errorHandlers: make(map[int]ErrorHandlerFunc),
		defaultHandler: func(c *gin.Context) {
			status := c.Writer.Status()
			c.JSON(status, gin.H{
				"error":       "An error occurred",
				"description": "No specific handler defined for this status",
			})
		},
	}
}

func (r *Router) SetDefaultErrorHandler(handler ErrorHandlerFunc) *Router {
	r.defaultHandler = handler
	r.router.NoRoute(func(cc *gin.Context) {
		handler(cc)
	})

	return r
}

// GetRoutes return list of routes
func (r *Router) GetRoutes() map[string]*Route {
	return r.routes
}

// SetErrorHandler set Error handler for status code
// Example:
//
//	router.SetErrorHandler(400, func(c *gin.Context) {
//			status := c.Writer.Status()
//			c.JSON(status, gin.H{
//				"error":       "An error occurred",
//				"description": "No specific handler defined for this status",
//			})
//		})
func (r *Router) SetErrorHandler(status int, handler ErrorHandlerFunc) *Router {

	if status == 404 {
		r.router.NoRoute(func(cc *gin.Context) {
			handler(cc)
		})
	}
	r.errorHandlers[status] = handler

	return r
}

// ErrorHandlerMiddleware middleware for customization error stats responses
func (r *Router) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()

		if status >= 400 {
			if handler, exists := r.errorHandlers[status]; exists {
				handler(c)
			} else {
				r.defaultHandler(c)
			}
		}
	}
}

// AddRouteList add RouteList to router
// Example:
//
//	 lr := NewRouteList()
//		v1 := CreateRouteList("/v1")
//		lr.addChild(v1)
//
//		lr.Add("/:locale/products/:id", productDetailHandler, Get)
//		v1.Add("/:locale/products/:id", productDetailHandler, Get)
//
//		router := NewRouter()
//		router.AddRouteList(lr)
func (r *Router) AddRouteList(l *RouteList) *Router {
	var group *gin.RouterGroup
	if l.pattern != "" {
		group = r.router.Group(l.pattern)
	}

	// Přidej všechny routy do skupiny nebo root routeru
	for _, route := range l.routes {
		if group != nil {
			createNativeRoute(*group, route)
			r.routes[route.name] = route
		} else {
			r.AddRoute(route.name, route.pattern, route.handler, route.method)
		}
	}

	// Rekurzivně přidej děti
	if l.children != nil { // Ověříme, že ukazatel na děti není nil
		for _, child := range l.children {
			r.AddRouteList(child)
		}
	}

	return r
}

// createHandlerFunc internal, add route to group, and return gin.IRoutes
func createNativeRoute(g gin.RouterGroup, route *Route) gin.IRoutes {
	return g.Handle(route.method.String(), route.pattern, createHandlerFunc(route.handler))
}

// createHandlerFunc internal, create gin.HandlerFunc
func createHandlerFunc(handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Získáme typ parametru druhého argumentu handleru
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 2 {

			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			return
		}

		// Vytvoříme nový prázdný objekt pro data
		paramType := handlerType.In(1) // Druhý parametr handleru (T)
		if paramType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("Handler parameter must be a pointer to a struct, got %v", paramType.Kind()))
		}
		paramElemType := paramType.Elem()
		paramValue := reflect.New(paramElemType).Interface()

		// Bindujeme data z požadavku do struktury
		err := c.ShouldBindUri(paramValue)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Zavoláme handler s bindenou strukturou
		reflect.ValueOf(handler).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(paramValue),
		})
	}
}

// AddRoute add route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
// method: Method
//
// Example:
//
//	router.AddRoute("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	}, Get)
func (r *Router) AddRoute(name string, pattern string, handler interface{}, method Method) *Router {
	r.routes[name] = NewRoute(name, pattern, handler, method, map[string]*Route{})

	r.router.Handle(method.String(), pattern, createHandlerFunc(handler))

	return r
}

// AddMultiMethodsRoute add route with multiple methods to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
// methods: []Method
//
// Example:
//
//	router.AddMultiMethodsRoute("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//			ID int `uri:"id" binding:"required"`
//		}) {
//		    // code
//			c.JSON(http.StatusOK, gin.H{
//				"id": p.ID,
//			})
//		}, []Method{Get, Post},
//		)
func (r *Router) AddMultiMethodsRoute(name string, pattern string, handler interface{}, methods []Method) *Router {
	for _, m := range methods {
		r.AddRoute(name, pattern, handler, m)
	}
	return r
}

// AddRouteMethod add route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
// method: Method
//
// Example:
//
//	router.AddRouteMethod("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	}, Get)
func (r *Router) AddRouteMethod(name string, pattern string, handler interface{}, method Method) *Router {
	return r.AddRoute(name, pattern, handler, method)
}

// CreateRoute add generic route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
// method: Method
//
// Example:
//
//	 CreateRoute(router, "user:detail", "/users/:id", func(c *gin.Context, p *struct{
//			ID int `uri:"id" binding:"required"`
//		}) {
//		    // code
//			c.JSON(http.StatusOK, gin.H{
//				"id": p.ID,
//			})
//		}, Get)
func CreateRoute[T any](r *Router, name string, pattern string, handler func(c *gin.Context, p *T), method Method) *Router {
	return r.AddRoute(name, pattern, handler, method)
}

// AddRouteGet add GET route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
//
// Example:
//
//	router.AddRouteGet("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	})
func (r *Router) AddRouteGet(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Get)
}

// AddRoutePost add POST route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
//
// Example:
//
//	router.AddRoutePost("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	})
func (r *Router) AddRoutePost(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Post)
}

// AddRoutePatch add Path route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
//
// Example:
//
//	router.AddRoutePatch("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	})
func (r *Router) AddRoutePatch(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Patch)
}

// AddRouteDelete add Delete route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
//
// Example:
//
//	router.AddRouteDelete("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	})
func (r *Router) AddRouteDelete(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Delete)
}

// AddRoutePut add Put route to router.
// name: name id for url construct
// pattern: path to route, example "/users/:id".
// handler: func(c *gin.Context, p *struct)
//
// Example:
//
//	router.AddRoutePut("/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//	    // code
//		c.JSON(http.StatusOK, gin.H{
//			"id": p.ID,
//		})
//	})
func (r *Router) AddRoutePut(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Put)
}

func (r *Router) AddRouteHead(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Head)
}

func (r *Router) AddRouteOptions(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Options)
}

func (r *Router) AddRouteConnect(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Connect)
}

func (r *Router) AddRouteTrace(name string, pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(name, pattern, handler, Trace)
}

// GetNativeRouter return gin router engine
// Docs continue in gin.Engine
func (r *Router) GetNativeRouter() *gin.Engine {
	return r.router
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (r *Router) Run(addr string) {
	router := r.router
	if h, ok := r.errorHandlers[404]; ok {
		router.NoRoute(func(c *gin.Context) {
			h(c)
		})
	} else {
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Custom 404",
			})
		})
	}
	router.Use(r.ErrorHandlerMiddleware())
	err := router.Run(addr)
	if err != nil {
		return
	}
}
