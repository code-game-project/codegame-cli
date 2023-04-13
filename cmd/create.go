package cmd

import (
	"github.com/spf13/cobra"

	"github.com/code-game-project/codegame-cli/create"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new CodeGame project",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		cobra.CheckErr(create.Create())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}