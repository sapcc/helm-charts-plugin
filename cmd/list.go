package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gosuri/uitable"
	"github.com/sapcc/helm-charts-plugin/pkg/charts"
	"github.com/sapcc/helm-outdated-dependencies/pkg/helm"
	"github.com/spf13/cobra"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

var listChartsLongUsage = `
Plugin to list Helm charts in the given folder. 

Examples:
  $ helm charts list <path> <flags>

  flags:
      --exclude-dirs strings   List of (sub-)directories to exclude.
      --include-vendor         Also consider charts in the vendor folder.
      --only-path              Only output the chart path.
      --output-dir string      If given, results will be written to file in this directory.
`

type listChartsCmd struct {
	helmSettings *helm_env.EnvSettings

	folder             string
	excludeDirs        []string
	timeout            time.Duration
	includeVendor      bool
	outputDir          string
	writeOnlyChartPath bool
}

func newListChartsCmd() *cobra.Command {
	l := &listChartsCmd{
		helmSettings: &helm_env.EnvSettings{
			Home: helm.GetHelmHome(),
		},
	}

	cmd := &cobra.Command{
		Use:          "list",
		Long:         listChartsLongUsage,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			folder, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			l.folder = folder

			if v, _ := cmd.Flags().GetStringSlice(flagExcludeDirs); v != nil {
				l.excludeDirs = v
			}

			if v, _ := cmd.Flags().GetString(flagOutputDir); v != "" {
				l.outputDir = v
			}

			if v, err := cmd.Flags().GetBool(flagWriteOnlyPath); err == nil {
				l.writeOnlyChartPath = v
			}

			if v, err := cmd.Flags().GetBool(flagIncludeVendor); err == nil {
				l.includeVendor = v
			}

			return l.list()
		},
	}

	addCommonFlags(cmd)

	return cmd
}

func (l *listChartsCmd) list() error {
	if !l.includeVendor {
		l.excludeDirs = append(l.excludeDirs, "vendor")
	}

	results, err := charts.ListHelmChartsInFolder(l.folder, l.excludeDirs)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("Not a single chart was found.")
		return nil
	}

	fmt.Println(l.formatTableOutput(results))

	if l.outputDir != "" {
		return l.writeToFile(results)
	}

	return nil
}

func (l *listChartsCmd) formatTableOutput(results []*charts.HelmChart) string {
	table := uitable.New()
	table.MaxColWidth = 200

	if !l.writeOnlyChartPath {
		table.AddRow("The following charts were found:")
		table.AddRow("NAME", "VERSION", "PATH")
	}

	for _, r := range results {
		if l.writeOnlyChartPath {
			table.AddRow(r.Path)
		} else {
			table.AddRow(r.Name, r.Version, r.Path)
		}
	}
	return table.String()
}

func (l *listChartsCmd) writeToFile(results []*charts.HelmChart) error {
	f, err := charts.EnsureFileExists(l.outputDir, outFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(l.formatTableOutput(results)))
	return err
}
