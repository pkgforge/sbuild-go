package linter

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type ValidationLevel struct {
	current int
	path    []string
	logger  *log.Logger
}

func (v *ValidationLevel) LogError(checkname string, err error) {
	fmt.Fprintf(os.Stderr, "[✗] %s: %v\n", checkname, err)
}

func (v *ValidationLevel) LogWarn(checkname string, warn string) {
	fmt.Fprintf(os.Stderr, "[-] %s: %v\n", checkname, warn)
}

func (v *ValidationLevel) LogSuccess(checkname, msg string) {
	fmt.Fprintf(os.Stderr, "[+] %s succeeded\n", checkname)
}

func (v *ValidationLevel) LogInfo(msg string) {
	fmt.Fprintf(os.Stderr, "[i] %s\n", msg)
}
