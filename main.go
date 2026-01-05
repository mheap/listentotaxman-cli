package main

import (
	"fmt"
	"os"

	"github.com/mheap/listentotaxman-cli/cmd"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	cmd.SetVersionInfo(Version, GitCommit, BuildDate)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
