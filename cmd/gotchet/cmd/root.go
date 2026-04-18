package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Root = &cobra.Command{
		Use:   "gotchet",
		Short: "Tool for reporting on Go test results",
	}
)
