package charts

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"k8s.io/helm/pkg/chartutil"
)

const chartMetadataName = "Chart.yaml"

// HelmChart is used to report the results of below functions.
type HelmChart struct {
	Name    string
	Version *semver.Version
	Path    string
}

// Equal checks if the given charts are equal.
func (h *HelmChart) Equal(c *HelmChart) bool {
	return h.Name == c.Name && h.Version.Equal(c.Version) && h.Path == c.Path
}

// ListHelmChartsInFolder list all Helm charts in the given folder.
func ListHelmChartsInFolder(folder string, excludeDirs []string, isUseRelativePath bool) ([]*HelmChart, error) {
	folder, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}

	var charts []*HelmChart
	err = filepath.Walk(folder, func(absPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && isValidChartDirectory(absPath, excludeDirs) {
			c, err := loadChartMetadata(absPath)
			if err != nil {
				return err
			}

			if isUseRelativePath {
				relPath, err := filepath.Rel(folder, c.Path)
				if err != nil {
					return err
				}
				c.Path = relPath
			}

			if !containsChart(charts, c) {
				charts = append(charts, c)
			}
		}
		return nil
	})

	return sortChartsAlphabetically(charts), err
}

// ListChangedHelmChartsInFolder compares the current version agains the given remote/branch:commit and lists the changed Helm charts.
func ListChangedHelmChartsInFolder(rootDirectory string, excludeDirs []string, remote, branch, commit string, isUseRelativePath bool) ([]*HelmChart, error) {
	git, err := newGit(rootDirectory, remote)
	if err != nil {
		return nil, err
	}

	if err := git.fetch(); err != nil {
		return nil, err
	}

	changedDirs, err := git.getChangedDirs(fmt.Sprintf("%s/%s", remote, branch), commit)
	if err != nil {
		return nil, err
	}

	var res []*HelmChart
	for _, dir := range changedDirs {
		path, err := getChartRootDirectory(rootDirectory, dir, excludeDirs)
		if err != nil {
			continue
		}

		c, err := loadChartMetadata(path)
		if err != nil {
			continue
		}

		if isUseRelativePath {
			relPath, err := filepath.Rel(rootDirectory, c.Path)
			if err != nil {
				continue
			}
			c.Path = relPath
		}

		if !containsChart(res, c) {
			res = append(res, c)
		}
	}
	return sortChartsAlphabetically(res), nil
}

func loadChartMetadata(absPathChartFolder string) (*HelmChart, error) {
	meta, err := chartutil.LoadChartfile(path.Join(absPathChartFolder, chartMetadataName))
	if err != nil {
		return nil, err
	}

	version, err := semver.NewVersion(meta.GetVersion())
	if err != nil {
		return nil, err
	}

	return &HelmChart{
		Name:    meta.GetName(),
		Version: version,
		Path:    absPathChartFolder,
	}, nil
}

func isValidChartDirectory(absPath string, excludeDirs []string) bool {
	if !filepath.IsAbs(absPath) {
		return false
	}

	for _, e := range excludeDirs {
		if strings.Contains(absPath, e) {
			return false
		}
	}

	_, err := os.Stat(path.Join(absPath, chartMetadataName))
	return !os.IsNotExist(err)
}

func getChartRootDirectory(root, path string, excludedDirs []string) (string, error) {
	if path == root {
		return "", errors.New("no more parent directories")
	}

	if isValidChartDirectory(path, excludedDirs) {
		return path, nil
	}

	return getChartRootDirectory(root, filepath.Dir(path), excludedDirs)
}

// containsCharts is used to avoid duplicates in a list of HelmCharts.
func containsChart(charts []*HelmChart, chart *HelmChart) bool {
	for _, c := range charts {
		if c.Equal(chart) {
			return true
		}
	}
	return false
}

func sortChartsAlphabetically(charts []*HelmChart) []*HelmChart {
	sort.Slice(charts, func(i, j int) bool {
		return charts[i].Name < charts[j].Name
	})
	return charts
}
