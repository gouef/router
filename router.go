package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

func (r *Router) GetRoutes() []interface{} {
	return r.routes
}

func (r *Router) GetRoutes2() gin.RoutesInfo {
	return r.routes2
}

type Router struct {
	router  *gin.Engine
	routes  []interface{}
	routes2 gin.RoutesInfo
}

func NewRouter() *Router {
	router := gin.Default()
	return &Router{router: router}
}

// Funkce pro přidání generické routy s generickým handlerem
func AddGenericRoute[T any](r *Router, pattern string, handler Handler[T], method Method) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 2 || handlerType.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic(fmt.Sprintf("Handler must be a function with signature func(*gin.Context, *T). Got: %s", handlerType))
	}

	// Přidáme handler do routeru
	r.routes = append(r.routes, handler)

	// Umožníme použít handler na základě typu
	r.router.Handle(method.String(), pattern, func(c *gin.Context) {
		dto, err := BindStruct[T](c) // Bindování dat do struktury
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Zavoláme handler s konkrétním typem
		handler(c, &dto)
	})
}

// V této metodě nebudeme dělat type assertion, použijeme generiku
func (r *Router) AddRoute(pattern string, handler interface{}, method Method) *Router {
	r.routes = append(r.routes, handler)

	// Předání handleru do Gin routeru s bindováním
	r.router.Handle(method.String(), pattern, func(c *gin.Context) {
		// Získáme typ parametru druhého argumentu handleru
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 2 {

			reflect.ValueOf(handler).Call([]reflect.Value{
				reflect.ValueOf(c),
			})
			//panic("Handler must be a function with signature func(*gin.Context, *T)")
			return
		}

		// Vytvoříme nový prázdný objekt pro data
		paramType := handlerType.In(1).Elem() // Druhý parametr handleru (T)
		paramValue := reflect.New(paramType).Interface()

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
	})
	r.routes2 = r.router.Routes()

	return r
}

// Přidání více metod pro stejnou routu
func (r *Router) AddMultiMethodsRoute(pattern string, handler interface{}, methods []Method) *Router {
	for _, m := range methods {
		r.AddRoute(pattern, handler, m)
	}
	return r
}

// Příklad metod pro různé HTTP metody
func (r *Router) AddRouteGet(pattern string, handler interface{}) *Router {
	return r.AddRoute(pattern, handler.(Handler[any]), Get)
}

func (r *Router) AddRoutePost(pattern string, handler interface{}) *Router {
	return r.AddRoute(pattern, handler.(Handler[any]), Post)
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

/*
// AddRoute přidá routu pro libovolný handler, který má (c *gin.Context, p *T)
func (r *Router) AddRoute(path string, handler interface{}, method Method) {
	// Ověříme, že handler je funkce s parametry (*gin.Context, *T)
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}
	if handlerType.NumIn() != 2 {
		panic("handler must have exactly two input parameters")
	}
	if handlerType.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic("first parameter must be *gin.Context")
	}

	h := func(c *gin.Context) {
		// Dynamicky vytvoříme prázdnou instanci parametru
		paramType := handlerType.In(1)
		paramInstance := reflect.New(paramType).Interface()

		// Bindujeme parametry URI
		if err := c.ShouldBindUri(paramInstance); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Zavoláme handler s parametry
		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(paramInstance)})
	}
	// Dynamicky přidáme route pro GET
	r.router.GET(path, h)

	// Zaznamenáme tuto routu pro testování
	r.routes = append(r.routes, gin.RouteInfo{
		Method:      method.String(),
		Path:        path,
		HandlerFunc: h,
	})
}
*/
/*
func (r *Router) AddRoute(path string, handler func(c *gin.Context, p interface{}), method Method, dto any) *Router {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	// Zajištění, že handler je funkce se dvěma argumenty: *gin.Context a nějaká struktura
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 2 || handlerType.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic(fmt.Sprintf("Handler must be a function with signature func(*gin.Context, T). Got: %s", handlerType))
	}
	//dtoType := handlerType.In(1)
	AddGenericRoute(r, r.router, pattern, handler.(Handler[any]), method)
	//AddGenericRoute(r, r.router, path, handler, method)

	return r
}

func (r *Router) AddMultiMethodsRoute(pattern string, handler Handler[T], methods []Method) *Router {
	for _, m := range methods {
		r.AddRoute(pattern, handler, m)
	}
	return r
}

func (r *Router) AddRouteGet(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Get)
}

func (r *Router) AddRoutePost(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Post)
}

func (r *Router) AddRoutePatch(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Patch)
}

func (r *Router) AddRouteDelete(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Delete)
}

func (r *Router) AddRoutePut(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Put)
}

func (r *Router) AddRouteHead(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Head)
}

func (r *Router) AddRouteOptions(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Options)
}

func (r *Router) AddRouteConnect(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Connect)
}

func (r *Router) AddRouteTrace(pattern string, handler Handler[T]) *Router {
	return r.AddRoute(pattern, handler, Trace)
}

func (r *Router) GetNativeRouter() *gin.Engine {
	return r.router
}

func (r *Router) Run(addr string) {
	//r.router.Use(RouteValidator(r.routes))
	err := r.router.Run(addr)
	if err != nil {
		return
	}
}
*/
