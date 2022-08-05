package cmd

import (
	"github.com/code-game-project/go-utils/cgfile"
	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "A CLI interface for share.code-game.org.",
}

func findGameURL() string {
	projectRoot, err := cgfile.FindProjectRoot()
	if err != nil {
		return ""
	}

	config, err := cgfile.LoadCodeGameFile(projectRoot)
	if err != nil {
		return ""
	}

	return config.URL
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
