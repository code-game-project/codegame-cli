package cmd

import (
	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "A CLI interface for share.code-game.org.",
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
