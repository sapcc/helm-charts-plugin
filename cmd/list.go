// Copyright 2025 SAP SE
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	helm_env "k8s.io/helm/pkg/helm/environment"

	"github.com/sapcc/helm-charts-plugin/pkg/charts"
)

var listChartsLongUsage = `
Plugin to list Helm charts in the given folder.

Examples:
  $ helm charts list <path> <flags>

  flags:
      --exclude-dirs        strings     List of (sub-)directories to exclude.
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
	useRelativePath,
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

			useRelativePath, err := cmd.Flags().GetBool(flagUseRelativePath)
			if err != nil {
				return err
			}
			l.useRelativePath = useRelativePath

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
	results, err := charts.ListHelmChartsInFolder(l.folder, l.excludeDirs, l.useRelativePath)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return errors.New("not a single chart was found")
	}

	table := FormatTableOutput(results, "The following charts were found:", l.writeOnlyChartPath, l.writeOnlyChartName)
	fmt.Println(table)

	if l.outputDir != "" {
		return l.writeToFile(table)
	}

	return nil
}

func (l *listChartsCmd) writeToFile(table string) error {
	f, err := charts.EnsureFileExists(l.outputDir, l.outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(table))
	return err
}

func FormatTableOutput(results []*charts.HelmChart, header string, writeOnlyChartPath, writeOnlyChartName bool) string {
	table := uitable.New()
	table.MaxColWidth = 200

	if !writeOnlyChartPath && !writeOnlyChartName {
		table.AddRow(header)
		table.AddRow("NAME", "VERSION", "PATH")
	}

	for _, r := range results {
		switch {
		case writeOnlyChartPath:
			table.AddRow(r.Path)
		case writeOnlyChartName:
			table.AddRow(r.Name)
		default:
			table.AddRow(r.Name, r.Version, r.Path)
		}
	}
	return table.String()
}
