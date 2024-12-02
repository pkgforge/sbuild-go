package linter

import (
	"fmt"
	"strings"
)

func (v *Validator) validateRequiredFields() ([]byte, error, string) {
	required := []string{"pkg", "description", "src_url", "x_exec.shell", "x_exec.run"}

	for _, field := range required {
		parts := strings.Split(field, ".")
		current := v.data

		for i, part := range parts {
			if val, ok := current[part]; !ok {
				return nil, fmt.Errorf("required field missing: %s", field), ""
			} else if i == len(parts)-1 {
				if val == nil || val == "" {
					return nil, fmt.Errorf("required field empty: %s", field), ""
				}
			} else if mapVal, ok := val.(map[string]interface{}); ok {
				current = mapVal
			} else {
				return nil, fmt.Errorf("invalid structure for field: %s", field), ""
			}
		}
	}

	// Check for pkgver or x_exec.pkgver
	if _, pkgverExists := v.data["pkgver"]; !pkgverExists {
		if xExec, ok := v.data["x_exec"].(map[string]interface{}); ok {
			if _, pkgverExists := xExec["pkgver"]; !pkgverExists {
				return nil, fmt.Errorf("pkgver field and x_exec.pkgver script are both missing"), ""
			}
		} else {
			return nil, fmt.Errorf("pkgver field and x_exec.pkgver script are both missing"), ""
		}
	}

	return nil, nil, ""
}
