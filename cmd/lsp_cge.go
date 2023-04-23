package cmd

import (
	"github.com/spf13/cobra"

	"github.com/code-game-project/codegame-cli/lsp/cge"
)

// lspCgeCmd represents the lsp cge command
var lspCgeCmd = &cobra.Command{
	Use:   "cge",
	Short: "Launch cge-ls",
	Run: func(_ *cobra.Command, _ []string) {
		err := cge.RunLSP()
		checkErr("Failed to launch cge-ls: %s", err)
	},
}

func init() {
	lspCmd.AddCommand(lspCgeCmd)
}
