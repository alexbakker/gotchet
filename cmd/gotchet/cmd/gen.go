package cmd

import (
	"fmt"
	"os"

	"github.com/alexbakker/gotchet/internal/report"
	"github.com/spf13/cobra"
)

var (
	genCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate a Go test report",
		Run:   startGen,
	}
)

func init() {
	genCmd.Flags().BoolVar(&rootFlags.Emulate, "emulate", false, "emulate run time of the test report (useful for development)")
	Root.AddCommand(genCmd)
}

func startGen(cmd *cobra.Command, args []string) {
	if err := report.Render(capture, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render report: %v", err)
	}
}
