package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/code-game-project/cli-utils/feedback"

	"github.com/code-game-project/codegame-cli/version"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "codegame-cli",
	Short: "The official CodeGame CLI",
}

func checkErr(format string, err error) {
	if err == nil {
		return
	}
	feedback.Fatal("codegame-cli", format, err)
	os.Exit(1)
}

func Execute() {
	rootCmd.SetVersionTemplate("codegame-cli {{.Version}}\n")
	rootCmd.Version = version.Version
	rootCmd.InitDefaultVersionFlag()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	feedback.Enable(feedback.NewCLIFeedback(feedback.SeverityInfo))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
