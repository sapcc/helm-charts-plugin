package main

import (
	"os"

	"github.com/sapcc/helm-charts-plugin/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
