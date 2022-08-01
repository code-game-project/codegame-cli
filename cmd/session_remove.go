package cmd

import (
	"fmt"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// sessionRemoveCmd represents the session remove command
var sessionRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a session.",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var gameURL string
		var username string
		var err error
		if len(args) > 0 {
			gameURL = args[0]
			if len(args) > 1 {
				username = args[1]
			}
		}

		if gameURL == "" {
			urls, err := sessions.ListGames()
			if err != nil {
				gameURL, err = cli.Input("Game URL:")
			} else {
				var index int
				index, err = cli.Select("Game URL:", urls)
				gameURL = urls[index]
			}
			if err != nil {
				return
			}
		}

		if username == "" {
			usernames, err := sessions.ListUsernames(gameURL)
			if err != nil {
				username, err = cli.Input("Username:")
			} else {
				var index int
				index, err = cli.Select("Username:", usernames)
				username = usernames[index]
			}
			if err != nil {
				return
			}
		}

		yes, err := cli.YesNo(fmt.Sprintf("Are you sure you want to remove '%s/%s'?", gameURL, username), false)
		if !yes {
			if err == nil {
				cli.Error("Canceled.")
			}
			return
		}

		err = sessions.Session{
			GameURL:  gameURL,
			Username: username,
		}.Remove()
		if err != nil {
			cli.Error("Failed to remove session: %s", err)
			return
		}
		cli.Success("Successfully removed session.")
	},
}

func init() {
	sessionCmd.AddCommand(sessionRemoveCmd)
}
