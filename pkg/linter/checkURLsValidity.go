package linter

import (
	"fmt"
	"strings"
)

func (v *Validator) validateURLs() ([]byte, error, string) {
	urlFields := map[string]bool{
		"desktop":     true,
		"icon":        true,
		"src_url":     false,
		"homepage":    false,
		"build_asset": false,
	}

	for field, singleValue := range urlFields {
		urls, exists := v.data[field]
		if !exists {
			continue
		}

		if singleValue {
			urlStr, ok := urls.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be a string", field), ""
			}

			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
				return nil, fmt.Errorf("%s contains invalid URL: %s", field, urlStr), ""
			}
		} else {
			urlSlice, ok := urls.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%s must be an array", field), ""
			}

			for _, url := range urlSlice {
				if field == "build_asset" {
					urlMap, ok := url.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("invalid URL type in %s", field), ""
					}

					if urlStr, ok := urlMap["url"].(string); ok {
						if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
							return nil, fmt.Errorf("%s contains invalid URL: %s", field, urlStr), ""
						}
					} else {
						return nil, fmt.Errorf("missing 'url' field in %s", field), ""
					}

					if _, ok := urlMap["out"].(string); !ok {
						return nil, fmt.Errorf("missing 'out' field in %s", field), ""
					}
				} else {
					urlStr, ok := url.(string)
					if !ok {
						return nil, fmt.Errorf("invalid URL type in %s", field), ""
					}

					if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
						return nil, fmt.Errorf("%s contains invalid URL: %s", field, urlStr), ""
					}
				}
			}
		}
	}

	return nil, nil, ""
}
