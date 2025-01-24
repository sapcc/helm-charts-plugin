package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gosuri/uitable"
	"github.com/sapcc/helm-charts-plugin/pkg/charts"
	"github.com/spf13/cobra"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

var listChartsLongUsage = `
Plugin to list Helm charts in the given folder.

Examples:
  $ helm charts list <path> <flags>

  flags:
      --exclude-dirs        strings     List of (sub-)directories to exclude.
      --include-vendor      bool        Also consider charts in the vendor folder.
      --only-path           bool        Only output the chart path.
      --output-dir          string      If given, results will be written to file in this directory.
      --output-filename     strin       Filename to use for output. (default "results.txt")
`

type listChartsCmd struct {
	helmSettings *helm_env.EnvSettings

	excludeDirs []string
	folder,
	outputDir,
	outputFilename string
	includeVendor,
	isUseRelativePath,
	writeOnlyChartPath,
	writeOnlyChartName bool
}

func newListChartsCmd() *cobra.Command {
	l := &listChartsCmd{
		helmSettings: &helm_env.EnvSettings{
			Home: charts.GetHelmHome(),
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

			excludeDirs, err := cmd.Flags().GetStringSlice(flagExcludeDirs)
			if err != nil {
				return err
			}
			if excludeDirs != nil {
				l.excludeDirs = excludeDirs
			}

			outputDir, err := cmd.Flags().GetString(flagOutputDir)
			if err != nil {
				return err
			}
			if outputDir != "" {
				l.outputDir = outputDir
			}

			outputFileName, err := cmd.Flags().GetString(flagOutputFileName)
			if err != nil {
				return err
			}
			if outputFileName != "" {
				l.outputFilename = outputFileName
			}

			includeVendor, err := cmd.Flags().GetBool(flagIncludeVendor)
			if err != nil {
				return err
			}
			l.includeVendor = includeVendor

			useRelativePath, err := cmd.Flags().GetBool(flagUseRelativePath)
			if err != nil {
				return err
			}
			l.isUseRelativePath = useRelativePath

			writeOnlyPath, err := cmd.Flags().GetBool(flagWriteOnlyPath)
			if err != nil {
				return err
			}
			l.writeOnlyChartPath = writeOnlyPath

			writeOnlyName, err := cmd.Flags().GetBool(flagWriteOnlyName)
			if err != nil {
				return err
			}
			l.writeOnlyChartName = writeOnlyName

			return l.list()
		},
	}

	addCommonFlags(cmd)

	return cmd
}

func (l *listChartsCmd) list() error {
	if !l.includeVendor {
		l.excludeDirs = append(l.excludeDirs, excludeVendorPaths...)
	}

	results, err := charts.ListHelmChartsInFolder(l.folder, l.excludeDirs, l.isUseRelativePath)
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

	if !l.writeOnlyChartPath && !l.writeOnlyChartName {
		table.AddRow("The following charts were found:")
		table.AddRow("NAME", "VERSION", "PATH")
	}

	for _, r := range results {
		if l.writeOnlyChartPath {
			table.AddRow(r.Path)
		} else if l.writeOnlyChartName {
			table.AddRow(r.Name)
		} else {
			table.AddRow(r.Name, r.Version, r.Path)
		}
	}
	return table.String()
}

func (l *listChartsCmd) writeToFile(results []*charts.HelmChart) error {
	f, err := charts.EnsureFileExists(l.outputDir, l.outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(l.formatTableOutput(results)))
	return err
}
