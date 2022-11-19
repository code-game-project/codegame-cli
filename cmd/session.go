package cmd

import (
	"errors"
	"os"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage CodeGame sessions.",
}

func selectSession(args []string) (sessions.Session, error) {
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
		if len(urls) == 0 {
			abort(errors.New("no sessions available"))
		}
		if err != nil {
			gameURL, err = cli.Input("Game URL:")
		} else {
			var index int
			index, err = cli.Select("Game URL:", urls)
			if err == nil {
				gameURL = urls[index]
			}
		}
		if err != nil {
			os.Exit(0)
		}
	}

	if username == "" {
		usernames, err := sessions.ListUsernames(gameURL)
		if len(usernames) == 0 {
			abort(errors.New("no sessions available"))
		}
		if err != nil {
			username, err = cli.Input("Username:")
		} else {
			var index int
			index, err = cli.Select("Username:", usernames)
			if err == nil {
				username = usernames[index]
			}
		}
		if err != nil {
			os.Exit(0)
		}
	}

	return sessions.LoadSession(gameURL, username)
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
