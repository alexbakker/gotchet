package main

import (
	"os"

	"github.com/alexbakker/gotchet/cmd/gotchet/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}
