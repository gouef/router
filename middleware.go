package router

/*
func RouteValidator(routes []Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Najdeme odpovídající routu
		for _, route := range routes {
			if route.Path == c.FullPath() {
				// Pokud routa odpovídá, provedeme validaci parametrů
				handlerType := reflect.TypeOf(route.Handler)

				// Validujeme pouze, pokud má handler parametr typu DTO
				if handlerType.Kind() == reflect.Func && handlerType.NumIn() == 2 {
					dtoType := handlerType.In(1) // Druhý parametr handleru

					// Kontrola požadovaných parametrů na základě DTO
					if dtoType.Kind() == reflect.Struct {
						validateDTO(c, dtoType)
					}
				}
				break
			}
		}

		// Pokud validace nevrátila chybu, pokračujeme
		c.Next()
	}
}

// Funkce pro validaci DTO
func validateDTO(c *gin.Context, dtoType reflect.Type) {
	for i := 0; i < dtoType.NumField(); i++ {
		field := dtoType.Field(i)
		paramName := field.Tag.Get("param")
		if paramName != "" {
			// Získání hodnoty parametru
			paramValue := c.Param(paramName)
			if paramValue == "" {
				// Pokud je hodnota prázdná, vrátíme chybu
				c.AbortWithStatusJSON(400, gin.H{
					"error": fmt.Sprintf("missing required parameter: %s", paramName),
				})
				return
			}

			// Další validace podle typu pole (např. čísla, stringy, atd.)
			if err := validateFieldType(paramValue, field.Type); err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"error": fmt.Sprintf("invalid value for parameter %s: %v", paramName, err),
				})
				return
			}
		}
	}
}

// Validace typu pole
func validateFieldType(value string, fieldType reflect.Type) error {
	switch fieldType.Kind() {
	case reflect.String:
		// String nepotřebuje extra validaci
		return nil
	case reflect.Int, reflect.Int64:
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("expected integer, got: %s", value)
		}
	case reflect.Float32, reflect.Float64:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("expected float, got: %s", value)
		}
	default:
		return fmt.Errorf("unsupported type: %s", fieldType.Kind())
	}
	return nil
}
*/
