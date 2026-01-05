package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of listentotaxman CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("listentotaxman version %s\n", version)
		if gitCommit != "unknown" {
			fmt.Printf("  commit: %s\n", gitCommit)
		}
		if buildDate != "unknown" {
			fmt.Printf("  built: %s\n", buildDate)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
