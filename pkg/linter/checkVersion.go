package linter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (v *Validator) validatePkgverScript(noShellcheck bool) ([]byte, error, string) {
	var version string
	var source string
	var warning string

	if xExec, ok := v.data["x_exec"].(map[string]interface{}); ok {
		shell, _ := xExec["shell"].(string)

		if script, ok := xExec["pkgver"].(string); ok {
			if strings.TrimSpace(script) == "" {
				return nil, fmt.Errorf("pkgver script is empty"), ""
			}

			validShell, warning := "", ""
			var err error
			if noShellcheck {
				warning = "shellcheck for x_exec.pkgver skipped"
			} else {
				validShell, err, warning = v.validateScript(script, shell)
				if err != nil {
					return nil, err, warning
				}
			}

			// Execute pkgver script
			cmd := exec.Command(validShell, "-c", script)
			output, err := cmd.CombinedOutput()
			if err != nil {
				exitErr := err.(*exec.ExitError)
				if exitErr.ExitCode() != 0 {
					return nil, fmt.Errorf("pkgver script exited with non-zero status: %d", exitErr.ExitCode()), warning
				}
				return nil, fmt.Errorf("failed to execute pkgver script: %w", err), warning
			}

			version = strings.TrimSpace(string(output))
			source = ".x_exec.pkgver"

			if version == "" {
				return nil, fmt.Errorf("pkgver script returned an empty version"), warning
			}

			// Write output to pkgver file
			if err := os.WriteFile(v.file+".pkgver", output, 0644); err != nil {
				return nil, fmt.Errorf("failed to write pkgver file: %w", err), warning
			}
		}
	}

	if version == "" {
		if pkgver, ok := v.data["pkgver"].(string); ok && pkgver != "" {
			version = pkgver
			source = ".pkgver"
		} else {
			return nil, fmt.Errorf("pkgver field is missing or empty"), warning
		}
	}

	if version != "" {
		v.level.LogInfo(fmt.Sprintf("Fetched Version: %s from [%s]", version, source))
	}

	return nil, nil, warning
}
