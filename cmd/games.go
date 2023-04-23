package cmd

import (
	"github.com/spf13/cobra"
)

// gamesCmd represents the games command
var gamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Manage games on a game server",
}

func init() {
	rootCmd.AddCommand(gamesCmd)
}
