package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
)

var Types = []string{
	"text",
	"number",
	"date",
	"boolean",
	"json",
	"file",
	"image",
	"video",
	"richtext",
}

func CheckTypes(data []byte) (bool, error) {
	var def map[string]interface{}
	if err := json.Unmarshal(data, &def); err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}

	typeSet := make(map[string]struct{}, len(Types))
	for _, t := range Types {
		typeSet[t] = struct{}{}
	}

	for key, val := range def {
		switch v := val.(type) {
		case string:
			if _, ok := typeSet[v]; !ok {
				return false, fmt.Errorf("invalid type for %q: %s", key, v)
			}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					if _, found := typeSet[s]; !found {
						return false, fmt.Errorf("invalid type in array for %q: %s", key, s)
					}
				} else {
					return false, fmt.Errorf("non-string value in array for %q", key)
				}
			}
		default:
			return false, fmt.Errorf("invalid type for %q: %T", key, v)
		}
	}
	return true, nil
}

func CompareSchemaWithData(schemaDef []byte, data map[string]interface{}) (bool, error) {
	var schema map[string]interface{}
	if err := json.Unmarshal(schemaDef, &schema); err != nil {
		return false, fmt.Errorf("invalid schema JSON: %w", err)
	}

	for key, schemaVal := range schema {
		val, exists := data[key]
		if !exists {
			return false, fmt.Errorf("missing required field %q", key)
		}

		if err := matchType(schemaVal, val); err != nil {
			return false, fmt.Errorf("field %q: %w", key, err)
		}
	}

	for key := range data {
		if _, exists := schema[key]; !exists {
			return false, fmt.Errorf("field %q not defined in schema", key)
		}
	}

	return true, nil
}

func matchType(schemaVal interface{}, value interface{}) error {
	switch t := schemaVal.(type) {
	case string:
		if !isPrimitiveTypeMatching(t, value) {
			return fmt.Errorf("expected %s, got %T", t, value)
		}
	case []interface{}:
		arr, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
		if len(t) == 0 {
			return nil
		}
		elemTypeStr, ok := t[0].(string)
		if !ok {
			return fmt.Errorf("invalid schema array type: %v", t[0])
		}
		for i, item := range arr {
			if !isPrimitiveTypeMatching(elemTypeStr, item) {
				return fmt.Errorf("element %d: expected %s, got %T", i, elemTypeStr, item)
			}
		}
	default:
		return fmt.Errorf("unsupported schema type %T", t)
	}
	return nil
}

func isPrimitiveTypeMatching(expectedType string, value interface{}) bool {
	switch expectedType {
	case "text", "string", "richtext":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "json":
		_, ok := value.(map[string]interface{})
		return ok
	case "file":
		_, ok := value.(string)
		return ok
	case "image", "video":
		str, ok := value.(string)
		if !ok {
			return false
		}
		// Validate URL format for media types
		u, err := url.ParseRequestURI(str)
		return err == nil && u.Scheme != "" && u.Host != ""
	default:
		return false
	}
}
