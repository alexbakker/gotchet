package cmd

import (
	"fmt"

	"github.com/alexbakker/gotchet/internal/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version information",
		Run:   startVersion,
	}
)

func init() {
	Root.AddCommand(versionCmd)
}

func startVersion(cmd *cobra.Command, args []string) {
	vs, err := version.String()
	if err != nil {
		exitWithError(err.Error())
		return
	}

	fmt.Print(vs)
	if ts := version.HumanRevisionTime(); ts != "" {
		fmt.Printf(" (%s)", ts)
	}
	fmt.Println()
	fmt.Println("https://github.com/alexbakker/gotchet (GPLv3)")
}
