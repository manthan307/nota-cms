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

type Field struct {
	Name       string      `json:"name"`
	Type       interface{} `json:"type"` // can be string or []string
	IsRequired bool        `json:"isRequired"`
}

// Validate schema definition syntax
func CheckTypes(data []byte) (bool, error) {
	var fields []Field
	if err := json.Unmarshal(data, &fields); err != nil {
		return false, fmt.Errorf("invalid JSON array: %w", err)
	}

	if len(fields) == 0 {
		return false, fmt.Errorf("schema must contain at least one field")
	}

	typeSet := make(map[string]struct{}, len(Types))
	for _, t := range Types {
		typeSet[t] = struct{}{}
	}

	for _, f := range fields {
		if f.Name == "" {
			return false, fmt.Errorf("each field must have a 'name'")
		}

		if f.Type == nil {
			return false, fmt.Errorf("field %q missing 'type'", f.Name)
		}

		switch t := f.Type.(type) {
		case string:
			if _, ok := typeSet[t]; !ok {
				return false, fmt.Errorf("field %q: invalid type %q", f.Name, t)
			}
		case []interface{}:
			for _, item := range t {
				s, ok := item.(string)
				if !ok {
					return false, fmt.Errorf("field %q: non-string type in array", f.Name)
				}
				if _, found := typeSet[s]; !found {
					return false, fmt.Errorf("field %q: invalid type in array %q", f.Name, s)
				}
			}
		default:
			return false, fmt.Errorf("field %q: type must be string or array", f.Name)
		}
	}

	return true, nil
}

// Compare schema definition with actual data
func CompareSchemaWithData(schemaDef []byte, data map[string]interface{}) (bool, error) {
	var fields []Field
	if err := json.Unmarshal(schemaDef, &fields); err != nil {
		return false, fmt.Errorf("invalid schema JSON array: %w", err)
	}

	fieldMap := make(map[string]Field)
	for _, f := range fields {
		fieldMap[f.Name] = f
	}

	for _, f := range fields {
		val, exists := data[f.Name]

		if !exists {
			if f.IsRequired {
				return false, fmt.Errorf("missing required field %q", f.Name)
			}
			continue
		}

		if err := matchType(f.Type, val); err != nil {
			return false, fmt.Errorf("field %q: %w", f.Name, err)
		}
	}

	for key := range data {
		if _, exists := fieldMap[key]; !exists {
			return false, fmt.Errorf("field %q not defined in schema", key)
		}
	}

	return true, nil
}

// Match single field's type
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

// Type matching logic
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
		u, err := url.ParseRequestURI(str)
		return err == nil && u.Scheme != "" && u.Host != ""
	default:
		return false
	}
}
