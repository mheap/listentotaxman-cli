package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
	buildDate string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "listentotaxman",
	Short: "Calculate UK tax and national insurance",
	Long: `listentotaxman is a CLI tool for calculating UK tax and national insurance.
It uses the listentotaxman.com API to provide accurate tax calculations.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information
func SetVersionInfo(v, commit, date string) {
	version = v
	gitCommit = commit
	buildDate = date
}

func init() {
	// Add subcommands here
}
