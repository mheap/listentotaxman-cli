package cmd

import (
	"os"

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
	// Add completion command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(listentotaxman completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ listentotaxman completion bash > /etc/bash_completion.d/listentotaxman
  # macOS:
  $ listentotaxman completion bash > /usr/local/etc/bash_completion.d/listentotaxman

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ listentotaxman completion zsh > "${fpath[1]}/_listentotaxman"
  # You will need to start a new shell for this setup to take effect.

fish:
  $ listentotaxman completion fish | source
  # To load completions for each session, execute once:
  $ listentotaxman completion fish > ~/.config/fish/completions/listentotaxman.fish

PowerShell:
  PS> listentotaxman completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> listentotaxman completion powershell > listentotaxman.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	})
}
