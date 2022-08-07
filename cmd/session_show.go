package cmd

import (
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// sessionShowCmd represents the session show command
var sessionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the session data.",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		session, err := selectSession(args)
		abortf("Failed to load session: %s", err)

		out := colorable.NewColorableStdout()
		printInfoProperty(out, "Game URL", session.GameURL, 14)
		printInfoProperty(out, "Username", session.Username, 14)
		printInfoProperty(out, "Game ID", session.GameId, 14)
		printInfoProperty(out, "Player ID", session.PlayerId, 14)
		printInfoProperty(out, "Player Secret", session.PlayerSecret, 14)
	},
}

func init() {
	sessionCmd.AddCommand(sessionShowCmd)
}
