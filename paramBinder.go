package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

type Handler[T any] func(c *gin.Context, p *T)

func BindStruct[T any](c *gin.Context) (T, error) {
	var dto T
	val := reflect.ValueOf(&dto).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		paramName := field.Tag.Get("param")

		if paramName != "" {
			paramValue := c.Param(paramName)
			fieldValue := val.Field(i)

			if fieldValue.CanSet() {
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
