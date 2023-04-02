package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/code-game-project/codegame-cli/docs"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: fmt.Sprintf("Open %s in the default web browser", docs.DocsURL),
	Run: func(_ *cobra.Command, _ []string) {
		err := docs.Open()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
