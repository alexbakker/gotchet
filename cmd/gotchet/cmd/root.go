package cmd

import (
	"fmt"
	"os"

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
)

func init() {
	Root.Flags().StringVarP(&rootFlags.Input, "input", "i", "-", "input filename (or - for stdin)")
	Root.Flags().BoolVarP(&rootFlags.Emulate, "emulate", "e", false, "emulate run time of the test report (useful for development)")
}

func startRoot(cmd *cobra.Command, args []string) {
	capture := runCapture(rootFlags.Input, rootFlags.Emulate)

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
