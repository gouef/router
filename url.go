package router

import (
	"fmt"
	"strings"
)

// GenerateUrlByPattern generate url by pattern
func GenerateUrlByPattern(pattern string, params map[string]interface{}) (string, error) {
	var urlBuilder strings.Builder
	isFirst := true

	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {

			paramName := part[1:] // Remove ':'
			if value, exists := params[paramName]; exists {
				if !isFirst {
					urlBuilder.WriteString("/")
				}

				urlBuilder.WriteString(fmt.Sprintf("%v", value))
				isFirst = false
			} else {
				return "", fmt.Errorf("missing value for parameter: %s", paramName)
			}
		} else {
			if !isFirst {
				urlBuilder.WriteString("/")
			}
			urlBuilder.WriteString(part)
			isFirst = false
		}
	}

	return urlBuilder.String(), nil
}
