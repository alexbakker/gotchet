package cmd

import (
	"fmt"
	"os"
)

func exitWithError(s string) {
	fmt.Fprintf(os.Stderr, s)
	os.Exit(1)
}
