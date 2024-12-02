package linter

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed embedded/validCategories.list
var categoriesTXT []byte

func (v *Validator) validateCategories() ([]byte, error, string) {
	var warn string

	// Check if the category field exists, set a default if not.
	categoryField, exists := v.data["category"]
	if !exists || categoryField == nil {
		v.data["category"] = []string{"Utility"}
		warn = fmt.Sprintf("adding default category: %s", v.data["category"])
		categoryField = v.data["category"]
	}

	// Parse allowed categories from the embedded file.
	allowedCategories := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(string(categoriesTXT)))
	for scanner.Scan() {
		category := scanner.Text()
		if category != "" {
			allowedCategories[category] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read categories: %w", err), ""
	}

	// Extract user-provided categories.
	var categories []string
	switch v := categoryField.(type) {
	case string:
		// Split by comma for multiple categories.
		categories = strings.Split(v, ",")
	case []interface{}:
		// Convert array of interfaces to array of strings.
		for _, cat := range v {
			if catStr, ok := cat.(string); ok {
				categories = append(categories, catStr)
			} else {
				return nil, fmt.Errorf("invalid category type in array"), ""
			}
		}
	case []string:
		categories = v
	case nil:
		categories = []string{"Utility"}
	default:
		return nil, fmt.Errorf("unsupported category format"), ""
	}

	// Validate user-provided categories against the allowed categories.
	for _, cat := range categories {
		if _, exists := allowedCategories[cat]; !exists {
			return nil, fmt.Errorf("invalid category: %s", cat), ""
		}
	}

	// Apply the default category if needed and issue a warning.
	if warn != "" {
		data, err := v.updateField([]string{"category"}, "Utility")
		if err != nil {
			return nil, fmt.Errorf("failed to update field: %w", err), ""
		}
		return data, nil, warn
	}

	return nil, nil, ""
}
