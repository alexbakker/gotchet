package cmd

import (
	"fmt"
	"os"

	"github.com/alexbakker/gotchet/internal/format"
)

func exitWithError(s string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", s)
	os.Exit(1)
}

func runCapture() *format.TestCapture {
	r := os.Stdin
	if rootFlags.Input != "-" {
		file, err := os.Open(rootFlags.Input)
		if err != nil {
			exitWithError(fmt.Sprintf("failed to open input: %v", err))
			return nil
		}
		defer file.Close()
		r = file
	}

	capture, err := format.Read(r, rootFlags.Emulate)
	if err != nil {
		exitWithError(fmt.Sprintf("failed to read test output: %v", err))
		return nil
	}

	return capture
}
