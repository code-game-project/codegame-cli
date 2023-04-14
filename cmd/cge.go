package cmd

import (
	"github.com/spf13/cobra"

	"github.com/code-game-project/codegame-cli/lsp/cge"
)

// cgeCmd represents the cge command
var cgeCmd = &cobra.Command{
	Use:   "cge",
	Short: "Launch cge-ls",
	Run: func(_ *cobra.Command, _ []string) {
		err := cge.RunLSP()
		checkErr("Failed to launch cge-ls: %s", err)
	},
}

func init() {
	lspCmd.AddCommand(cgeCmd)
}
