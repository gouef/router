package router

import (
	"fmt"
	"strings"
)

// GenerateUrlByPattern generate url by Pattern
func GenerateUrlByPattern(pattern string, params map[string]interface{}) (string, error) {
	if pattern == "" {
		return "/", nil
	}

	parts := strings.Split(pattern, "/")
	var resultParts []string

	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, ":") {
			paramName := part[1:]

			value, exists := params[paramName]
			if !exists || fmt.Sprintf("%v", value) == "" {
				return "", fmt.Errorf("missing value for parameter: %s", paramName)
			}

			resultParts = append(resultParts, fmt.Sprintf("%v", value))
		} else {
			resultParts = append(resultParts, part)
		}
	}

	return "/" + strings.Join(resultParts, "/"), nil
}
