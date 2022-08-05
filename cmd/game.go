package cmd

import (
	"github.com/spf13/cobra"
)

// gameCmd represents the game command
var gameCmd = &cobra.Command{
	Use:   "game",
	Short: "Manage CodeGame games.",
}

func init() {
	rootCmd.AddCommand(gameCmd)
}
