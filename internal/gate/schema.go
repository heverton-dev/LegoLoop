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
	return matchType(schema, val) && matchRequired(schema, val) && matchProperties(schema, val, depth) && matchItems(schema, val, depth)
}

// matchType confere "type". Tipo ausente, que não é texto ou desconhecido não restringe.
func matchType(schema map[string]any, val any) bool {
	typeStr, ok := schema["type"].(string)
	if !ok {
		return true
	}
	switch typeStr {
	case "object":
		_, ok := val.(map[string]any)
		return ok
	case "array":
		_, ok := val.([]any)
		return ok
	case "string":
		_, ok := val.(string)
		return ok
	case "number":
		_, ok := val.(float64)
		return ok
	case "integer":
		f, ok := val.(float64)
		return ok && math.Trunc(f) == f
	case "boolean":
		_, ok := val.(bool)
		return ok
	case "null":
		return val == nil
	}
	return true
}

// matchRequired confere "required": o valor é objeto e tem todas as chaves, que devem ser texto.
func matchRequired(schema map[string]any, val any) bool {
	reqList, ok := schema["required"].([]any)
	if !ok {
		return true
	}
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
	return true
}

// matchProperties confere "properties": o valor é objeto e cada propriedade presente cumpre o seu sub-esquema.
func matchProperties(schema map[string]any, val any, depth int) bool {
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return true
	}
	obj, ok := val.(map[string]any)
	if !ok {
		return false
	}
	for propName, propDef := range props {
		fieldVal, exists := obj[propName]
		if !exists {
			continue
		}
		if subSchema, ok := propDef.(map[string]any); ok && !matchSchema(subSchema, fieldVal, depth+1) {
			return false
		}
	}
	return true
}

// matchItems confere "items": o valor é lista e cada item cumpre o sub-esquema.
func matchItems(schema map[string]any, val any, depth int) bool {
	itemsDef, ok := schema["items"].(map[string]any)
	if !ok {
		return true
	}
	arr, ok := val.([]any)
	if !ok {
		return false
	}
	for _, item := range arr {
		if !matchSchema(itemsDef, item, depth+1) {
			return false
		}
	}
	return true
}
