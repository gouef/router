package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gouef/mode"
	"net/http"
	"reflect"
)

const (
	// DebugMode indicates Mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates Mode is release.
	ReleaseMode = "release"
	// TestMode indicates Mode is test.
	TestMode = "test"
)

type ErrorHandlerFunc func(c *gin.Context)

type Router struct {
	Router         *gin.Engine
	Routes         map[string]*Route
	middlewares    []interface{}
	ErrorHandlers  map[int]ErrorHandlerFunc
	DefaultHandler ErrorHandlerFunc
	Mode           *mode.Mode
}

type routeGroupStruct struct {
	pattern  string
	routes   []*Route
	children []*routeGroupStruct
}

// NewRouter create new Router
func NewRouter() *Router {
	router := gin.New()
	router.Use(GouefMiddleware())
	m, _ := mode.NewBasicMode()
	return &Router{
		Router:        router,
		Routes:        make(map[string]*Route),
		ErrorHandlers: make(map[int]ErrorHandlerFunc),
		DefaultHandler: func(c *gin.Context) {
			status := c.Writer.Status()
			c.JSON(status, gin.H{
				"error":       "An error occurred",
				"description": "No specific Handler defined for this status",
			})
		},
		Mode: m,
	}
}

// SetDefaultErrorHandler set default error Handler
// Example:
//
//	router.SetDefaultErrorHandler(func(c *gin.Context) {
//		c.JSON(http.StatusNotFound, gin.H{
//		"error": "Custom 404",
//	})
func (r *Router) SetDefaultErrorHandler(handler ErrorHandlerFunc) *Router {
	r.DefaultHandler = handler
	r.Router.NoRoute(func(cc *gin.Context) {
		handler(cc)
	})

	return r
}

// GetRoutes return list of Routes
func (r *Router) GetRoutes() map[string]*Route {
	return r.Routes
}

// SetErrorHandler set Error Handler for status code
// Example:
//
//	router.SetErrorHandler(400, func(c *gin.Context) {
//		status := c.Writer.Status()
//		c.JSON(status, gin.H{
//			"error":       "An error occurred",
//			"description": "No specific Handler defined for this status",
//		})
//	})
func (r *Router) SetErrorHandler(status int, handler ErrorHandlerFunc) *Router {

	if status == 404 {
		r.Router.NoRoute(func(cc *gin.Context) {
			handler(cc)
		})
	}
	r.ErrorHandlers[status] = handler

	return r
}

// ErrorHandlerMiddleware middleware for customization error stats responses
func (r *Router) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()

		if status >= 400 {
			if handler, exists := r.ErrorHandlers[status]; exists {
				handler(c)
			} else {
				r.DefaultHandler(c)
			}
		}
	}
}

// AddRouteList add RouteList to router
// Example:
//
//	lr := NewRouteList()
//	v1 := CreateRouteList("/v1")
//	lr.addChild(v1)
//
//	lr.Add("/:locale/products/:id", productDetailHandler, Get)
//	v1.Add("/:locale/products/:id", productDetailHandler, Get)
//
//	router := NewRouter()
//	router.AddRouteList(lr)
func (r *Router) AddRouteList(l *RouteList) error {
	groupStruct := r.createGroupStruct(l)
	return r.createNativeGroup(groupStruct)
}

// createNativeGroup internal, create gin groups
func (r *Router) createNativeGroup(rgs *routeGroupStruct) error {
	var group *gin.RouterGroup
	group = r.GetNativeRouter().Group(rgs.pattern)

	for _, route := range rgs.routes {
		_, err := createNativeRoute(*group, route)
		r.Routes[route.Name] = route
		if err != nil {
			return err
		}
	}

	return r.createNativeGroupChildren(group, rgs.children)
}

// createNativeGroupChildren internal, create gin groups children
func (r *Router) createNativeGroupChildren(group *gin.RouterGroup, children []*routeGroupStruct) error {

	for _, child := range children {
		childGroup := group.Group(child.pattern)
		for _, route := range child.routes {
			_, err := createNativeRoute(*childGroup, route)
			if err != nil {
				return err
			}

			r.Routes[route.Name] = route
		}
		if len(child.children) > 0 {
			err2 := r.createNativeGroupChildren(childGroup, child.children)

			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

// createGroupStruct internal, create structure for generate gin groups
func (r *Router) createGroupStruct(l *RouteList) *routeGroupStruct {
	s := &routeGroupStruct{
		pattern: l.pattern,
		routes:  l.routes,
	}

	var children []*routeGroupStruct

	if l.children != nil {
		for _, child := range l.children {
			children = append(children, r.createGroupStruct(child))
		}
	}
	s.children = children

	return s
}

// createHandlerFunc internal, add route to group, and return gin.IRoutes
func createNativeRoute(g gin.RouterGroup, route *Route) (gin.IRoutes, error) {
	handler, err := createHandlerFunc(route.Handler)
	if err != nil {
		return nil, err
	}
	return g.Handle(route.Method.String(), route.Pattern, handler), nil
}

// createHandlerFunc internal, create gin.HandlerFunc
func createHandlerFunc(handler interface{}) (gin.HandlerFunc, error) {
	handlerType := reflect.TypeOf(handler)

	switch handlerType {
	case reflect.TypeOf((*HandlerContext)(nil)).Elem():
		return func(c *gin.Context) {
			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			return
		}, nil
	/*case reflect.TypeOf((*HandlerContextParam)(nil)).Elem():
	return func(c *gin.Context) {
		paramType := handlerType.In(1)
		paramElemType := paramType.Elem()
		paramValue := reflect.New(paramElemType).Interface()

		err := c.ShouldBindUri(paramValue)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Call Handler with binded struct
		reflect.ValueOf(handler).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(paramValue),
		})
		return
	}, nil*/
	case reflect.TypeOf((*HandlerFunc)(nil)).Elem():
		return func(c *gin.Context) {
			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			return
		}, nil
	}

	if handlerType.Kind() != reflect.Func {
		return nil, errors.New("Handler must be a function")
	}

	numIn := handlerType.NumIn()
	if numIn < 1 || numIn > 2 {
		return nil, errors.New("Handler must have one or two parameters")
	}

	if numIn == 2 {
		paramType := handlerType.In(1)
		if paramType.Kind() != reflect.Ptr {
			return nil, errors.New(fmt.Sprintf("Handler parameter must be a pointer to a struct, got %v", paramType.Kind()))
		}
	}

	return func(c *gin.Context) {
		if numIn == 1 {
			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			return
		} else {
			paramType := handlerType.In(1)
			paramElemType := paramType.Elem()
			paramValue := reflect.New(paramElemType).Interface()

			err := c.ShouldBindUri(paramValue)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Call Handler with binded struct
			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
				reflect.ValueOf(paramValue),
			})
		}
	}, nil
}

// AddRoute add route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
// Method: Method
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
func (r *Router) AddRoute(name string, pattern string, handler interface{}, method Method) error {

	handlerFunc, err := createHandlerFunc(handler)
	if err != nil {
		return err
	}

	r.Routes[name] = NewRoute(name, pattern, handler, method)

	r.Router.Handle(method.String(), pattern, handlerFunc)

	return nil
}

func (r *Router) AddRouteObject(route *Route) *Router {
	l := NewRouteList()
	l.AddRoute(route)

	r.AddRouteList(l)
	return r
}

// AddMultiMethodsRoute add route with multiple methods to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
// methods: []Method
//
// Example:
//
//	router.AddMultiMethodsRoute("user:detail", "/users/:id", func(c *gin.Context, p *struct{
//		ID int `uri:"id" binding:"required"`
//	}) {
//		// code
//		c.JSON(http.StatusOK, gin.H{
//				"id": p.ID,
//			})
//		}, []Method{Get, Post},
//	)
func (r *Router) AddMultiMethodsRoute(name string, pattern string, handler interface{}, methods []Method) *Router {
	for _, m := range methods {
		r.AddRoute(name, pattern, handler, m)
	}
	return r
}

// AddRouteMethod add route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
// Method: Method
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
func (r *Router) AddRouteMethod(name string, pattern string, handler interface{}, method Method) error {
	return r.AddRoute(name, pattern, handler, method)
}

// CreateRoute add generic route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
// Method: Method
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
func CreateRoute[T any](r *Router, name string, pattern string, handler func(c *gin.Context, p *T), method Method) error {
	//res := &Response{Context: handler}
	return r.AddRoute(name, pattern, handler, method)
}

// AddRouteGet add GET route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
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
func (r *Router) AddRouteGet(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Get)
}

// AddRoutePost add POST route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
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
func (r *Router) AddRoutePost(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Post)
}

// AddRoutePatch add Path route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
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
func (r *Router) AddRoutePatch(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Patch)
}

// AddRouteDelete add Delete route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
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
func (r *Router) AddRouteDelete(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Delete)
}

// AddRoutePut add Put route to router.
// Name: Name id for url construct
// Pattern: path to route, example "/users/:id".
// Handler: func(c *gin.Context, p *struct)
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
func (r *Router) AddRoutePut(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Put)
}

func (r *Router) AddRouteHead(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Head)
}

func (r *Router) AddRouteOptions(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Options)
}

func (r *Router) AddRouteConnect(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Connect)
}

func (r *Router) AddRouteTrace(name string, pattern string, handler interface{}) error {
	return r.AddRouteMethod(name, pattern, handler, Trace)
}

func (r *Router) GenerateUrlByName(name string, params map[string]interface{}) (string, error) {
	route, exists := r.Routes[name]

	if !exists {
		return "", errors.New(fmt.Sprintf("route with Name %s not found", name))
	}

	return r.GenerateUrlByPattern(route.Pattern, params)
}

func (r *Router) GenerateUrlByPattern(pattern string, params map[string]interface{}) (string, error) {
	return GenerateUrlByPattern(pattern, params)
}

// SetHtmlRenderer set html renderer
func (r *Router) SetHtmlRenderer(renderer render.HTMLRender) *Router {
	r.GetNativeRouter().HTMLRender = renderer
	return r
}

func (r *Router) EnablePrefetch() *Router {
	r.GetNativeRouter().Use(func(c *gin.Context) {
		c.Writer.Header().Del("Purpose")
		c.Writer.Header().Set("Purpose", "prefetch")
		c.Writer.Header().Set("X-DNS-Prefetch-Control", "on")
		c.Next()
	})

	return r
}

// GetNativeRouter return gin router engine
// Docs continue in gin.Engine
func (r *Router) GetNativeRouter() *gin.Engine {
	return r.Router
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this Method will block the calling goroutine indefinitely unless an error happens.
func (r *Router) Run(addr string) error {
	router := r.Router
	if h, ok := r.ErrorHandlers[404]; ok {
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
	//router.Use(r.ErrorHandlerMiddleware())
	return router.Run(addr)
}
