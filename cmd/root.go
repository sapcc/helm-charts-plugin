package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagExcludeDirs   = "exclude-dirs"
	flagIncludeVendor = "include-vendor"
	flagOutputDir     = "output-dir"
	flagWriteOnlyPath = "only-path"
	outFileName       = "result.txt"
)

var rootCmdLongUsage = `
Plugin that helps to manage helm charts.

Examples:
  $ helm charts list 		 <path> <flags>		- List Helm charts in the given directory.
  $ helm charts list-changed <path> <flags> 	- Identify and list Helm charts that were changed compared to another commit.
`

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "charts",
		Long:      rootCmdLongUsage,
		ValidArgs: []string{"chartpath"},
	}

	cmd.AddCommand(
		newListChartsCmd(),
		newChangedChartsCmd(),
	)

	return cmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(flagExcludeDirs, "", []string{}, "List of (sub-)directories to exclude.")
	cmd.Flags().BoolP(flagIncludeVendor, "", false, "Also consider charts in the vendor folder.")
	cmd.Flags().StringP(flagOutputDir, "", "", "If given, results will be written to file in this directory.")
	cmd.Flags().BoolP(flagWriteOnlyPath, "", false, "Only output the chart path.")
}
