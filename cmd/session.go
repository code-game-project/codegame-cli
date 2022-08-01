package cmd

import (
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage CodeGame sessions.",
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
