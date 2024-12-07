package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

type Param struct {
}

type Handler[T any] func(c *gin.Context, p *T)

//type Handler func(c *gin.Context, p interface{})

func BindStruct[T any](c *gin.Context) (T, error) {
	var dto T
	val := reflect.ValueOf(&dto).Elem()

	// Pro každý field ve struktuře provedeme bind
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		paramName := field.Tag.Get("param") // Získá hodnotu z tagu `param`

		if paramName != "" {
			paramValue := c.Param(paramName)
			fieldValue := val.Field(i)

			if fieldValue.CanSet() {
				// Podle typu hodnoty provedeme správné nastavení
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(paramValue)
				case reflect.Int, reflect.Int32, reflect.Int64:
					intValue, err := strconv.Atoi(paramValue)
					if err != nil {
						return dto, fmt.Errorf("invalid integer value for %s: %v", paramName, err)
					}
					fieldValue.SetInt(int64(intValue))
				case reflect.Float32, reflect.Float64:
					floatValue, err := strconv.ParseFloat(paramValue, 64)
					if err != nil {
						return dto, fmt.Errorf("invalid float value for %s: %v", paramName, err)
					}
					fieldValue.SetFloat(floatValue)
				default:
					return dto, fmt.Errorf("unsupported field type %s", fieldValue.Kind())
				}
			}
		}
	}

	return dto, nil
}
