package router

import (
	"fmt"
	"strings"
)

func GenerateUrlByPattern(pattern string, params map[string]interface{}) (string, error) {
	// Předpokládáme, že pattern může obsahovat parametry jako :id
	var urlBuilder strings.Builder
	isFirst := true

	// Rozdělíme pattern podle / pro iteraci
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			// Pokud je část patternu parametr (např. :id), hledáme jeho hodnotu v params
			paramName := part[1:] // Odstraníme ':'
			if value, exists := params[paramName]; exists {
				if !isFirst {
					urlBuilder.WriteString("/")
				}
				// Přidáme hodnotu parametru do URL
				urlBuilder.WriteString(fmt.Sprintf("%v", value))
				isFirst = false
			} else {
				// Pokud parametr v params neexistuje, vrátíme chybu
				return "", fmt.Errorf("missing value for parameter: %s", paramName)
			}
		} else {
			// Pokud část patternu není parametr, přidáme ji přímo do URL
			if !isFirst {
				urlBuilder.WriteString("/")
			}
			urlBuilder.WriteString(part)
			isFirst = false
		}
	}

	return urlBuilder.String(), nil
}
