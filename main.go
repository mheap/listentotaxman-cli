// Package main is the entry point for the listentotaxman CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/mheap/listentotaxman-cli/cmd"
)

// Version information set at build time
var (
	// Version is the application version
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
