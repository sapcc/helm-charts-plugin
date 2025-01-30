package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gosuri/uitable"
	"github.com/sapcc/helm-charts-plugin/pkg/charts"
	"github.com/spf13/cobra"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

var changedChartsLongUsage = `
List Helm charts that were changed compared to a given Git commit.

Examples:
  $ helm charts list-changed <path> <flags>

  flags:
    --branch 			string			The name of the branch used to identify changes. (default "master")
    --commit 			string          The commit used to identify changes. (default "HEAD")
    --exclude-dirs 		strings   		List of (sub-)directories to exclude.
    --only-path         bool     		Only output the chart path.
    --output-dir 		string      	If given, results will be written to file in this directory.
    --output-filename 	string			Filename to use for output. (default "results.txt")
    --remote 			string          The name of the git remote used to identify changes. (default "origin)

`

type changedChartsCmd struct {
	helmSettings *helm_env.EnvSettings

	directory          string
	excludeDirs        []string
	outputDir          string
	outputFilename     string
	writeOnlyChartPath bool
	writeOnlyChartName bool
	isUseRelativePath  bool

	remote,
	branch,
	commit string
}

func newChangedChartsCmd() *cobra.Command {
	c := &changedChartsCmd{
		helmSettings: &helm_env.EnvSettings{
			Home: charts.GetHelmHome(),
		},
	}

	cmd := &cobra.Command{
		Use:          "list-changed",
		Long:         changedChartsLongUsage,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			d, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			c.directory = d

			excludeDirs, err := cmd.Flags().GetStringSlice(flagExcludeDirs)
			if err != nil {
				return err
			}
			c.excludeDirs = excludeDirs

			outputDir, err := cmd.Flags().GetString(flagOutputDir)
			if err != nil {
				return err
			}
			c.outputDir = outputDir

			outputFileName, err := cmd.Flags().GetString(flagOutputFileName)
			if err != nil {
				return err
			}
			c.outputFilename = outputFileName

			writeOnlyName, err := cmd.Flags().GetBool(flagWriteOnlyName)
			if err != nil {
				return err
			}
			c.writeOnlyChartName = writeOnlyName

			useRelativePath, err := cmd.Flags().GetBool(flagUseRelativePath)
			if err != nil {
				return err
			}
			c.isUseRelativePath = useRelativePath

			v, err := cmd.Flags().GetBool(flagWriteOnlyPath)
			if err != nil {
				return err
			}
			c.writeOnlyChartPath = v

			return c.listChanged()
		},
	}

	addCommonFlags(cmd)
	cmd.Flags().StringVarP(&c.remote, "remote", "", "origin", "The name of the git remote used to identify changes.")
	cmd.Flags().StringVarP(&c.branch, "branch", "", "master", "The name of the branch used to identify changes.")
	cmd.Flags().StringVarP(&c.commit, "commit", "", "HEAD", "The commit used to identify changes.")

	return cmd
}

func (c *changedChartsCmd) listChanged() error {
	results, err := charts.ListChangedHelmChartsInFolder(c.directory, c.excludeDirs, c.remote, c.branch, c.commit, c.isUseRelativePath)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("Nothing was changed.")
		return nil
	}

	fmt.Println(c.formatTableOutput(results))

	if c.outputDir != "" {
		return c.writeToFile(results)
	}

	return nil
}

func (c *changedChartsCmd) formatTableOutput(results []*charts.HelmChart) string {
	table := uitable.New()
	table.MaxColWidth = 200

	if !c.writeOnlyChartPath && !c.writeOnlyChartName {
		table.AddRow(fmt.Sprintf("Compared to %s/%s:%s following charts were changed:", c.remote, c.branch, c.commit))
		table.AddRow("NAME", "VERSION", "PATH")
	}

	for _, r := range results {
		if c.writeOnlyChartPath {
			table.AddRow(r.Path)
		} else if c.writeOnlyChartName {
			table.AddRow(r.Name)
		} else {
			table.AddRow(r.Name, r.Version, r.Path)
		}
	}
	return table.String()
}

func (c *changedChartsCmd) writeToFile(results []*charts.HelmChart) error {
	f, err := charts.EnsureFileExists(c.outputDir, c.outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(c.formatTableOutput(results)))
	return err
}
