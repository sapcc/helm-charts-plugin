package cmd

import "github.com/spf13/cobra"

const (
	flagExcludeDirs     = "exclude-dirs"
	flagIncludeVendor   = "include-vendor"
	flagOutputDir       = "output-dir"
	flagOutputFileName  = "output-filename"
	flagWriteOnlyPath   = "only-path"
	flagUseRelativePath = "relative-path"
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
	cmd.Flags().BoolP(flagIncludeVendor, "", false, "Whether to consider sub-charts in the /vendor and /charts folders.")
	cmd.Flags().StringP(flagOutputDir, "", "", "If given, results will be written to file in this directory.")
	cmd.Flags().StringP(flagOutputFileName, "", "results.txt", "Filename to use for output.")
	cmd.Flags().BoolP(flagWriteOnlyPath, "", false, "Only output the chart path.")
	cmd.Flags().BoolP(flagUseRelativePath, "", false, "Return chart path' relative to the given directory.")
}
