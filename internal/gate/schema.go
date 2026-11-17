package gate

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
)

// ValidateJSONSchema valida se targetBytes obedece às restrições definidas em schemaBytes.
// Implementa validação estrita sem dependências externas: type, required, properties, items,
// const, minimum e $ref local (#/$defs/...).
func ValidateJSONSchema(schemaBytes, targetBytes []byte) bool {
	var schema map[string]any
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return false
	}
	var target any
	if err := json.Unmarshal(targetBytes, &target); err != nil {
		return false
	}
	return matchSchema(schema, schema, target, 0)
}

// matchSchema confere o valor contra o sub-esquema. O esquema raiz resolve os $ref locais.
func matchSchema(root, schema map[string]any, val any, depth int) bool {
	if depth > MaxDepth {
		return false
	}
	if ref, ok := schema["$ref"].(string); ok {
		alvo, ok := resolverRef(root, ref)
		if !ok {
			return false
		}
		schema = alvo
	}
	return matchType(schema, val) && matchRequired(schema, val) && matchProperties(root, schema, val, depth) &&
		matchItems(root, schema, val, depth) && matchConst(schema, val) && matchMinimum(schema, val)
}

// resolverRef resolve um $ref local do tipo "#/$defs/<nome>" contra o esquema raiz.
func resolverRef(root map[string]any, ref string) (map[string]any, bool) {
	const prefixo = "#/$defs/"
	if !strings.HasPrefix(ref, prefixo) {
		return nil, false
	}
	defs, ok := root["$defs"].(map[string]any)
	if !ok {
		return nil, false
	}
	alvo, ok := defs[ref[len(prefixo):]].(map[string]any)
	return alvo, ok
}

// matchConst confere "const": o valor deve ser idêntico ao constante declarado.
func matchConst(schema map[string]any, val any) bool {
	alvo, existe := schema["const"]
	if !existe {
		return true
	}
	return reflect.DeepEqual(alvo, val)
}

// matchMinimum confere "minimum": o número deve ser maior ou igual ao limite.
func matchMinimum(schema map[string]any, val any) bool {
	limite, ok := schema["minimum"].(float64)
	if !ok {
		return true
	}
	n, ok := val.(float64)
	return ok && n >= limite
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
func matchProperties(root, schema map[string]any, val any, depth int) bool {
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
		if subSchema, ok := propDef.(map[string]any); ok && !matchSchema(root, subSchema, fieldVal, depth+1) {
			return false
		}
	}
	return true
}

// matchItems confere "items": o valor é lista e cada item cumpre o sub-esquema.
// Null não é lista: os itens só se aplicam a listas, como manda o contrato do JSON Schema.
func matchItems(root, schema map[string]any, val any, depth int) bool {
	itemsDef, ok := schema["items"].(map[string]any)
	if !ok {
		return true
	}
	if val == nil {
		return true
	}
	arr, ok := val.([]any)
	if !ok {
		return false
	}
	for _, item := range arr {
		if !matchSchema(root, itemsDef, item, depth+1) {
			return false
		}
	}
	return true
}
