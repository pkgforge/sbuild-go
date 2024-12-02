package linter

import (
	"fmt"
)

func (v *Validator) validateDisabledField() ([]byte, error, string) {
	if _, ok := v.data["_disabled"]; !ok {
		return nil, fmt.Errorf("_disabled field does not exist"), ""
	}

	if _, ok := v.data["_disabled"].(bool); !ok {
		return nil, fmt.Errorf("_disabled field must be a boolean"), ""
	}

	return nil, nil, ""
}
