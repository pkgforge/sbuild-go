package linter

import (
	"fmt"
	"strings"
)

func (v *Validator) validateNoEmptyFields() ([]byte, error, string) {
	var clean func(interface{}, string) interface{}
	var emptyFields []string
	var trailingSpaces []string

	clean = func(i interface{}, path string) interface{} {
		switch x := i.(type) {
		case map[string]interface{}:
			m2 := make(map[string]interface{})
			for k, v := range x {
				newPath := fmt.Sprintf("%s.%s", path, k)
				if v == nil || v == "" {
					emptyFields = append(emptyFields, newPath)
				} else {
					cleaned := clean(v, newPath)
					if cleaned != nil {
						m2[k] = cleaned
					}
				}
			}
			if len(m2) > 0 {
				return m2
			}
		case []interface{}:
			var s2 []interface{}
			for idx, v := range x {
				newPath := fmt.Sprintf("%s[%d]", path, idx)
				if v == nil || v == "" {
					emptyFields = append(emptyFields, newPath)
				} else {
					cleaned := clean(v, newPath)
					if cleaned != nil {
						s2 = append(s2, cleaned)
					}
				}
			}
			if len(s2) > 0 {
				return s2
			}
		case string:
			trimmed := strings.TrimSpace(x)
			if trimmed != x {
				trailingSpaces = append(trailingSpaces, path)
			}
			if trimmed != "" {
				return trimmed
			}
		default:
			if x != nil {
				return x
			}
		}
		return nil
	}

	cleaned := clean(v.data, "").(map[string]interface{})
	v.data = cleaned

	var warn string
	if len(emptyFields) > 0 {
		warn = fmt.Sprintf("Empty fields encountered: %v", emptyFields)
	}
	if len(trailingSpaces) > 0 {
		if warn != "" {
			warn += "; "
		}
		warn += fmt.Sprintf("Trailing spaces removed from: %v", trailingSpaces)
	}

	if warn != "" {
		data, err := v.updateField([]string{}, v.data)
		if err != nil {
			return nil, fmt.Errorf("failed to update field: %w", err), ""
		}
		return data, nil, warn
	}

	return nil, nil, ""
}
