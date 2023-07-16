package main

import (
	"fmt"
	"os"

	"github.com/alexbakker/gotchet/cmd/gotchet/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
