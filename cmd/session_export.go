package cmd

import (
	"github.com/spf13/cobra"
)

// sessionExportCmd represents the session export command
var sessionExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a session to share.code-game.org (same as 'codegame share session').",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		shareSessionCmd.Run(cmd, args)
	},
}

func init() {
	sessionCmd.AddCommand(sessionExportCmd)
}
