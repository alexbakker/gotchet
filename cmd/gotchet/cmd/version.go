package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	versionNumber       string
	versionRevision     string
	versionRevisionTime string
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
	if versionNumber == "" {
		exitWithError("no version information")
		return
	}

	fmt.Printf("v%s-%s", versionNumber, versionRevision)
	if ts := parseRevisionTime(versionRevisionTime); ts != "" {
		fmt.Printf(" (%s)", ts)
	}
	fmt.Println()
}

func parseRevisionTime(s string) string {
	secs, err := strconv.ParseInt(versionRevisionTime, 10, 64)
	if err != nil {
		return ""
	}

	return time.Unix(secs, 0).UTC().String()
}
