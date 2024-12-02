package linter

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	allowedPkgTypes = []string{
		"appbundle", "appimage", "archive", "dynamic",
		"flatimage", "gameimage", "nixappimage", "runimage", "static",
	}

	pkgPattern          = regexp.MustCompile(`^[a-zA-Z0-9\+\-_\.]+$`)
	invalidCharsPattern = regexp.MustCompile(`[^a-zA-Z0-9\+\-_\.]`)
)

func (v *Validator) validateEnforcedFields() ([]byte, error, string) {
	getInvalidChars := func(field string) string {
		matches := invalidCharsPattern.FindAllString(field, -1)
		if len(matches) > 0 {
			return strings.Join(matches, ", ")
		}
		return ""
	}

	if pkg, ok := v.data["pkg"].(string); ok {
		if invalidChars := getInvalidChars(pkg); invalidChars != "" {
			return nil, fmt.Errorf(".pkg contains invalid characters: %s", invalidChars), ""
		}
	}

	if pkgID, ok := v.data["app_id"].(string); ok {
		if invalidChars := getInvalidChars(pkgID); invalidChars != "" {
			return nil, fmt.Errorf(".app_id contains invalid characters: %s", invalidChars), ""
		}
	}

	if pkgType, ok := v.data["pkg_type"].(string); ok {
		valid := false
		for _, allowed := range allowedPkgTypes {
			if pkgType == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf(".pkg_type has invalid type: %s", pkgType), ""
		}
	}

	return nil, nil, ""
}
