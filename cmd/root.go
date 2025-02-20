package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "tempura",
	Short:   `A CLI tool that leverages templates to smoothly create projects.`,
	Long: `  ______________  _______  __  ____    ___ 
 /_  __/ ____/  |/  / __ \/ / / / /   /   |
  / / / __/ / /|_/ / /_/ / / / / /   / /| |
 / / / /___/ /  / / ____/ /_/ / /___/ ___ |
/_/ /_____/_/  /_/_/    \____/_____/_/  |_|

A CLI tool that leverages templates to smoothly create projects.

GitHub:
  https://github.com/nihi-lo/tempura`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
