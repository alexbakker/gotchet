package cmd

import (
	"fmt"
	"os"

	"github.com/alexbakker/gotchet/internal/format"
	"github.com/alexbakker/gotchet/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	Root = &cobra.Command{
		Use:   "gotchet",
		Short: "Go test reporter",
		Run:   startRoot,
	}
	rootFlags = struct {
		Input   string
		Emulate bool
	}{}
	capture *format.TestCapture
)

func init() {
	cobra.OnInitialize(startCapture)
	Root.PersistentFlags().StringVarP(&rootFlags.Input, "input", "i", "-", "input filename (or - for stdin)")
	Root.Flags().BoolVarP(&rootFlags.Emulate, "emulate", "e", false, "emulate run time of the test report (useful for development)")
}

func startCapture() {
	r := os.Stdin
	if rootFlags.Input != "-" {
		var err error
		file, err := os.Open(rootFlags.Input)
		if err != nil {
			exitWithError(fmt.Sprintf("Failed to open input: %v", err))
			return
		}
		defer file.Close()
		r = file
	}

	var err error
	capture, err = format.Read(r, rootFlags.Emulate)
	if err != nil {
		exitWithError(fmt.Sprintf("Failed to read test output: %v", err))
		return
	}
}

func startRoot(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(tui.New(capture),
		tea.WithMouseCellMotion(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "bye!")
}
