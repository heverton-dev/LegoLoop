package gate

import (
	"encoding/json"
	"math"
)

// ValidateJSONSchema valida se targetBytes obedece às restrições definidas em schemaBytes.
// Implementa validação estrita sem dependências externas: type, required, properties e items.
func ValidateJSONSchema(schemaBytes, targetBytes []byte) bool {
	var schema map[string]any
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return false
	}
	var target any
	if err := json.Unmarshal(targetBytes, &target); err != nil {
		return false
	}
	return matchSchema(schema, target, 0)
}

func matchSchema(schema map[string]any, val any, depth int) bool {
	if depth > MaxDepth {
		return false
	}

	// 1. Validação de type
	if typeStr, ok := schema["type"].(string); ok {
		switch typeStr {
		case "object":
			if _, ok := val.(map[string]any); !ok {
				return false
			}
		case "array":
			if _, ok := val.([]any); !ok {
				return false
			}
		case "string":
			if _, ok := val.(string); !ok {
				return false
			}
		case "number":
			if _, ok := val.(float64); !ok {
				return false
			}
		case "integer":
			f, ok := val.(float64)
			if !ok || math.Trunc(f) != f {
				return false
			}
		case "boolean":
			if _, ok := val.(bool); !ok {
				return false
			}
		case "null":
			if val != nil {
				return false
			}
		}
	}

	// 2. Validação de required
	if reqList, ok := schema["required"].([]any); ok {
		obj, ok := val.(map[string]any)
		if !ok {
			return false
		}
		for _, req := range reqList {
			key, ok := req.(string)
			if !ok {
				return false
			}
			if _, exists := obj[key]; !exists {
				return false
			}
		}
	}

	// 3. Validação de properties
	if props, ok := schema["properties"].(map[string]any); ok {
		obj, ok := val.(map[string]any)
		if !ok {
			return false
		}
		for propName, propDef := range props {
			fieldVal, exists := obj[propName]
			if !exists {
				continue
			}
			if subSchema, ok := propDef.(map[string]any); ok {
				if !matchSchema(subSchema, fieldVal, depth+1) {
					return false
				}
			}
		}
	}

	// 4. Validação de items
	if itemsDef, ok := schema["items"].(map[string]any); ok {
		arr, ok := val.([]any)
		if !ok {
			return false
		}
		for _, item := range arr {
			if !matchSchema(itemsDef, item, depth+1) {
				return false
			}
		}
	}

	return true
}
