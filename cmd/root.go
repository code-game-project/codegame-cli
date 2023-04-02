package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "codegame-cli",
	Short: "The official CodeGame CLI",
	Long:  "codegame-cli helps you develop CodeGame applications",
}

func Execute(version string) {
	rootCmd.SetVersionTemplate("codegame-cli {{.Version}}\n")
	rootCmd.Version = version
	rootCmd.InitDefaultVersionFlag()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
