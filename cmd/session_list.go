package cmd

import (
	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// sessionListCmd represents the session list command
var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available sessions.",
	Run: func(cmd *cobra.Command, args []string) {
		sessionList, err := sessions.ListSessions()
		if err != nil {
			cli.Error("Failed to retrieve session list: %s", err)
			return
		}
		for game, usernames := range sessionList {
			cli.PrintColor(cli.CyanBold, game)
			for _, u := range usernames {
				cli.Print("  - %s", u)
			}
		}
		if len(sessionList) == 0 {
			cli.Print("No sessions stored.")
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionListCmd)
}
