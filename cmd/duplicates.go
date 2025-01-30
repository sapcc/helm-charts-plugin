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

var findDuplicatesChartsLongUsage = `
Plugin to find duplicate Helm charts in the given folder.

Examples:
  $ helm charts find-duplicates <path> <flags>

  flags:
      --exclude-dirs				strings		  List of (sub-)directories to exclude.
      --only-path           bool   			Only output the chart path.
      --output-dir		    	string   		If given, results will be written to file in this directory.
      --output-filename     string   		Filename to use for output. (default "results.txt")
			--fail-on-duplicates	bool				Fail if duplicate charts are found.
`

type findDuplicatesChartsCmd struct {
	helmSettings *helm_env.EnvSettings
	folder,
	outputDir,
	outputFilename string
	writeOnlyChartPath,
	isUseRelativePath,
	failOnDuplicates bool
	excludeDirs []string
}

func newFindDuplicatesChartsCmd() *cobra.Command {
	l := &findDuplicatesChartsCmd{
		helmSettings: &helm_env.EnvSettings{
			Home: charts.GetHelmHome(),
		},
	}

	cmd := &cobra.Command{
		Use:          "find-duplicates",
		Long:         findDuplicatesChartsLongUsage,
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
			l.excludeDirs = excludeDirs

			outputDir, err := cmd.Flags().GetString(flagOutputDir)
			if err != nil {
				return err
			}
			l.outputDir = outputDir

			outputFileName, err := cmd.Flags().GetString(flagOutputFileName)
			if err != nil {
				return err
			}
			l.outputFilename = outputFileName

			writeOnlyPath, err := cmd.Flags().GetBool(flagWriteOnlyPath)
			if err != nil {
				return err
			}
			l.writeOnlyChartPath = writeOnlyPath

			useRelativePath, err := cmd.Flags().GetBool(flagUseRelativePath)
			if err != nil {
				return err
			}
			l.isUseRelativePath = useRelativePath

			return l.findDuplicates()
		},
	}

	addCommonFlags(cmd)
	cmd.Flags().BoolVarP(&l.failOnDuplicates, "fail-on-duplicates", "", false, "Fail if duplicate charts are found.")

	return cmd
}

func (l *findDuplicatesChartsCmd) findDuplicates() error {
	results, err := charts.FindDuplicateChartsInFolder(l.folder, l.excludeDirs, l.isUseRelativePath)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No duplicates found.")
		return nil
	}

	fmt.Println(l.formatTableOutput(results))

	if l.outputDir != "" {
		return l.writeToFile(results)
	}

	if l.failOnDuplicates {
		return errors.New("found multiple helm charts with the same name")
	}

	return nil
}

func (l *findDuplicatesChartsCmd) formatTableOutput(results []*charts.HelmChart) string {
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

func (l *findDuplicatesChartsCmd) writeToFile(results []*charts.HelmChart) error {
	f, err := charts.EnsureFileExists(l.outputDir, l.outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(l.formatTableOutput(results)))
	return err
}
