package cmd

import (
	"fmt"
	"os"

	"github.com/alexbakker/gotchet/internal/report"
	"github.com/spf13/cobra"
)

var (
	genCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate a Go test report",
		Run:     startGen,
	}
	genFlags = struct {
		Output string
	}{}
)

func init() {
	genCmd.Flags().StringVarP(&genFlags.Output, "output", "o", "-", "output filename (or - for stdout)")
	Root.AddCommand(genCmd)
}

func startGen(cmd *cobra.Command, args []string) {
	w := os.Stdout
	if genFlags.Output != "-" {
		file, err := os.Open(genFlags.Output)
		if err != nil {
			exitWithError(fmt.Sprintf("Failed to open output: %v", err))
			return
		}
		defer file.Close()
		w = file
	}

	if err := report.Render(capture, w); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render report: %v", err)
	}
}
