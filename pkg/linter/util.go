package linter

import (
	"fmt"
	"os"
	"os/exec"
)

// validateScript handles the common validation logic for shell scripts
func (v *Validator) validateScript(script string, shell string) (string, error, string) {
	var shebang string
	if shell != "" && shell != "sh" {
		shebang = fmt.Sprintf("#!/usr/bin/env %s", shell)
	} else {
		shell = "sh"
		shebang = "#!/bin/sh"
	}

	// Check if interpreter exists
	if _, err := exec.LookPath(shell); err != nil {
		return "", fmt.Errorf("interpreter (%s) not found in $PATH", shell), ""
	}

	fullScript := shebang + "\n" + script

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "shellcheck-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err), ""
	}
	defer os.Remove(tmpFile.Name())

	// Write script to temp file
	if _, err := tmpFile.WriteString(fullScript); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err), ""
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err), ""
	}

	cmdError := exec.Command("shellcheck", "--severity=error", tmpFile.Name())
	if output, err := cmdError.CombinedOutput(); err != nil {
		fmt.Println(fmt.Errorf("shellcheck with --severity=error failed: %s", string(output)))
		cmdWarning := exec.Command("shellcheck", "--severity=warning", tmpFile.Name())
		if output, err := cmdWarning.CombinedOutput(); err != nil {
			return "", nil, fmt.Sprintf("shellcheck warnings: %s", output)
		}
	}

	return shell, nil, ""
}

func removeDuplicates(data interface{}) interface{} {
	switch x := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range x {
			newMap[k] = removeDuplicates(v)
		}
		return newMap
	case []interface{}:
		seen := make(map[string]bool)
		var newList []interface{}
		for _, v := range x {
			strValue, ok := v.(string)
			if !ok || !seen[strValue] {
				newList = append(newList, v)
				seen[strValue] = true
			}
		}
		return newList
	default:
		return data
	}
}