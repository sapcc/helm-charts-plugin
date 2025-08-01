// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package charts

import (
	"fmt"
	"os"
	"path"

	"github.com/sapcc/go-bits/osext"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
)

// EnsureFileExists ensures all directories and the file itself exist.
func EnsureFileExists(absPath, filename string) (*os.File, error) {
	err := os.MkdirAll(absPath, os.ModeDir)
	if err != nil {
		return nil, err
	}

	filepath := path.Join(absPath, filename)
	fmt.Println("Using file: ", filepath)

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	err = f.Truncate(0)
	return f, err
}

// GetHelmHome returns the HELM_HOME path.
func GetHelmHome() helmpath.Home {
	return helmpath.Home(osext.GetenvOrDefault("HELM_HOME", helm_env.DefaultHelmHome))
}
