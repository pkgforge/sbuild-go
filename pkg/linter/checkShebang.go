package linter

import (
	"bufio"
	"bytes"
	"strings"
)

func (v *Validator) checkShebang() ([]byte, error, string) {
	scanner := bufio.NewScanner(bytes.NewReader(v.raw))
	var buffer bytes.Buffer

	// Scan through the lines to collect them
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteByte('\n')
	}

	// Check if the first line is the shebang
	if buffer.Len() > 0 && strings.HasPrefix(buffer.String(), "#!/SBUILD") {
		return nil, nil, ""
	}

	// Prepend the shebang to the collected lines
	newContent := bytes.NewBufferString("#!/SBUILD\n")
	newContent.Write(buffer.Bytes())

	// Update the raw field with the new content
	v.raw = newContent.Bytes()

	return newContent.Bytes(), nil, "didn't find a shebang in the SBUILD file, automatically fixed"
}
