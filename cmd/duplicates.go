package cmd

import (
	"errors"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/sapcc/helm-charts-plugin/pkg/charts"
	"github.com/sapcc/helm-outdated-dependencies/pkg/helm"
	"github.com/spf13/cobra"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"path/filepath"
)

var findDuplicatesChartsLongUsage = `
Plugin to find duplicate Helm charts in the given folder. 

Examples:
  $ helm charts find-duplicates <path> <flags>

  flags:
      --exclude-dirs				strings		  List of (sub-)directories to exclude.
      --include-vendor      bool   			Also consider charts in the vendor folder.
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
	includeVendor,
	writeOnlyChartPath,
	isUseRelativePath,
	failOnDuplicates bool
	excludeDirs []string
}

func newFindDuplicatesChartsCmd() *cobra.Command {
	l := &findDuplicatesChartsCmd{
		helmSettings: &helm_env.EnvSettings{
			Home: helm.GetHelmHome(),
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

			if v, _ := cmd.Flags().GetStringSlice(flagExcludeDirs); v != nil {
				l.excludeDirs = v
			}

			if v, _ := cmd.Flags().GetString(flagOutputDir); v != "" {
				l.outputDir = v
			}

			if v, _ := cmd.Flags().GetString(flagOutputFileName); v != "" {
				l.outputFilename = v
			}

			if v, err := cmd.Flags().GetBool(flagWriteOnlyPath); err == nil {
				l.writeOnlyChartPath = v
			}

			if v, err := cmd.Flags().GetBool(flagIncludeVendor); err == nil {
				l.includeVendor = v
			}

			if v, err := cmd.Flags().GetBool(flagUseRelativePath); err == nil {
				l.isUseRelativePath = v
			}

			return l.findDuplicates()
		},
	}

	addCommonFlags(cmd)
	cmd.Flags().BoolVarP(&l.failOnDuplicates, "fail-on-duplicates", "", false, "Fail if duplicate charts are found.")

	return cmd
}

func (l *findDuplicatesChartsCmd) findDuplicates() error {
	if !l.includeVendor {
		l.excludeDirs = append(l.excludeDirs, excludeVendorPaths...)
	}

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
