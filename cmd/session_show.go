package cmd

import (
	"fmt"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// sessionShowCmd represents the session show command
var sessionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the session data.",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var gameURL string
		var username string
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

		session, err := sessions.LoadSession(gameURL, username)
		if err != nil {
			abort(fmt.Errorf("The session %s@%s does not exist!", username, gameURL))
		}

		out := colorable.NewColorableStdout()
		printInfoProperty(out, "Name", fmt.Sprintf("%s@%s", username, gameURL), 14)
		printInfoProperty(out, "Game ID", session.GameId, 14)
		printInfoProperty(out, "Player ID", session.PlayerId, 14)
		printInfoProperty(out, "Player Secret", session.PlayerSecret, 14)
	},
}

func init() {
	sessionCmd.AddCommand(sessionShowCmd)
}
