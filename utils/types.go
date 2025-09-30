package utils

import (
	"encoding/json"
	"fmt"
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
	// "audio",
	"richtext",
}

func CheckTypes(data []byte) (bool, error) {
	// Unmarshal into generic map
	var def map[string]interface{}
	if err := json.Unmarshal(data, &def); err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}

	// Build a quick lookup map for Types
	typeSet := make(map[string]struct{}, len(Types))
	for _, t := range Types {
		typeSet[t] = struct{}{}
	}

	// Validate each field
	for key, val := range def {
		switch v := val.(type) {
		case string:
			if _, ok := typeSet[v]; !ok {
				return false, fmt.Errorf("invalid type for %q: %s", key, v)
			}
		case []interface{}: // if it's an array of strings
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
