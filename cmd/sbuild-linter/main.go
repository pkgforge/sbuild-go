package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

var Log *log.Logger

const errorMessage = `incorrect SBUILD File. Please recheck @ https://www.yamllint.com
  SBUILD docs: https://github.com/pkgforge/soarpkgs/blob/main/SBUILD.md
  SBUILD Specification: https://github.com/pkgforge/soarpkgs/blob/main/SBUILD_SPEC.md
`

func main() {
	Log = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})

	pkgverFlag := flag.Bool("pkgver", false, "Enable pkgver validation")
	noShellcheckFlag := flag.Bool("no-shellcheck", false, "Disables shellcheck usage in pkgver & run script validation")
	flag.Parse()

	if flag.NArg() < 1 {
		Log.Fatal("Usage: sbuild-validator <file1> [<file2> ...]")
	}

	warningCount := 0
	errorCount := 0
	successCount := 0

	// Filter out flags from the positional arguments
	files := make([]string, 0, flag.NArg())
	for _, arg := range flag.Args() {
		if arg == "--pkgver" {
			continue
		}
		files = append(files, arg)
	}

	for _, file := range files {
		// Print which file is being verified
		fmt.Printf("\x1b[44m\x1b[30m\x1b[4m[+]\x1b[0m Verifying %s\n", file)

		validator, err := NewValidator(file)
		if err != nil {
			Log.Error(err.Error())
			Log.Error(errorMessage)
			errorCount++
			continue
		}

		validatedData, warnings, err := validator.ValidateAll(*pkgverFlag, *noShellcheckFlag)
		if err != nil {
			Log.Error(errorMessage)
			errorCount++
			continue
		}

		// Handle warnings and success separately
		if warnings > 0 {
			fmt.Printf("[!] %s has %d warnings\n", file, warnings)
			warningCount++
		} else {
			fmt.Printf("[✓] %s is valid\n", file)
			successCount++
		}
		println()

		if err := writeDataToNewFile(file, validatedData); err != nil {
			Log.Error("Failed to write validated data", "file", file, "error", err)
			errorCount++
			continue
		}
		println()

	}

	// Print summary statistics
	fmt.Printf("Validation Summary:\n")
	fmt.Printf("Files with warnings: %d\n", warningCount)
	fmt.Printf("Files with errors: %d\n", errorCount)
	fmt.Printf("Files passing all checks: %d\n", successCount)

	// Exit with error if any files had errors
	if errorCount > 0 {
		os.Exit(1)
	}
}
