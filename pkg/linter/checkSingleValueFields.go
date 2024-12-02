package linter

import (
	"fmt"
)

func (v *Validator) validateSingleValueFields() ([]byte, error, string) {
	singleValueFields := []string{"pkg", "pkg_id", "pkg_type", "description"}

	for _, field := range singleValueFields {
		if val, ok := v.data[field]; ok {
			switch val.(type) {
			case string:
			case []interface{}:
				return nil, fmt.Errorf("field '%s' must be a single value, not an array. Value: %v", field, val), ""
			default:
				return nil, fmt.Errorf("field '%s' must be a string. Value: %v", field, val), ""
			}
		}
	}
	return nil, nil, ""
}
