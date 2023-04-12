package cmd

import (
	"github.com/spf13/cobra"
)

// lspCmd represents the lsp command
var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Launch a language server",
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
