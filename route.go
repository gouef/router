package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Route[T any] struct {
	Path    string
	Handler Handler[T]
	Dto     T
}

func CreateRoute[T any](pattern string, handler Handler[T]) *Route[T] {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	// Zajištění, že handler je funkce se dvěma argumenty: *gin.Context a nějaká struktura
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 2 || handlerType.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic(fmt.Sprintf("Handler must be a function with signature func(*gin.Context, T). Got: %s", handlerType))
	}

	//c := handlerType.In(0)

	//dto, err := BindStruct[T](&c)

	return &Route[T]{
		Path:    pattern,
		Handler: handler,
	}
}
