package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

func (r *Router) GetRoutes() []interface{} {
	return r.routes
}

type Router struct {
	router *gin.Engine
	routes []interface{}
}

func NewRouter() *Router {
	router := gin.Default()
	return &Router{router: router}
}

func (r *Router) AddRouteList(l *RouteList) *Router {
	var group *gin.RouterGroup
	if l.pattern != "" {
		group = r.router.Group(l.pattern)
	}

	// Přidej všechny routy do skupiny nebo root routeru
	for _, route := range l.routes {
		if group != nil {
			createNativeRoute(*group, route)
			r.routes = append(r.routes, route.handler)
		} else {
			r.AddRoute(route.pattern, route.handler, route.method)
			r.routes = append(r.routes, route.handler)
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

func createNativeRoute(g gin.RouterGroup, route Route) gin.IRoutes {
	return g.Handle(route.method.String(), route.pattern, createHandlerFunc(route.handler))
}

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

// V této metodě nebudeme dělat type assertion, použijeme generiku
func (r *Router) AddRoute(pattern string, handler interface{}, method Method) *Router {
	r.routes = append(r.routes, handler)

	// Předání handleru do Gin routeru s bindováním
	r.router.Handle(method.String(), pattern, createHandlerFunc(handler))

	return r
}

// Přidání více metod pro stejnou routu
func (r *Router) AddMultiMethodsRoute(pattern string, handler interface{}, methods []Method) *Router {
	for _, m := range methods {
		r.AddRoute(pattern, handler, m)
	}
	return r
}

func (r *Router) AddRouteMethod(pattern string, handler interface{}, method Method) *Router {
	return r.AddRoute(pattern, handler.(Handler[any]), method)
}

// Příklad metod pro různé HTTP metody
func (r *Router) AddRouteGet(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Get)
}

func (r *Router) AddRoutePost(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Post)
}

func (r *Router) AddRoutePatch(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Patch)
}

func (r *Router) AddRouteDelete(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Delete)
}

func (r *Router) AddRoutePut(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Put)
}

func (r *Router) AddRouteHead(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Head)
}

func (r *Router) AddRouteOptions(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Options)
}

func (r *Router) AddRouteConnect(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Connect)
}

func (r *Router) AddRouteTrace(pattern string, handler interface{}) *Router {
	return r.AddRouteMethod(pattern, handler, Trace)
}

// A tak dál pro další metody
func (r *Router) GetNativeRouter() *gin.Engine {
	return r.router
}

func (r *Router) Run(addr string) {
	err := r.router.Run(addr)
	if err != nil {
		return
	}
}
