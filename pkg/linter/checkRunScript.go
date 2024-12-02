package linter

func (v *Validator) validateRunScript(noShellcheck bool) ([]byte, error, string) {
	if xExec, ok := v.data["x_exec"].(map[string]interface{}); ok {
		shell, _ := xExec["shell"].(string)
		if script, ok := xExec["run"].(string); ok {

			if noShellcheck {
				return nil, nil, "shellcheck for x_exec.run skipped"
			}

			_, err, warning := v.validateScript(script, shell)
			if err != nil {
				return nil, err, warning
			}
			return nil, nil, warning
		}
	}
	return nil, nil, ""
}
