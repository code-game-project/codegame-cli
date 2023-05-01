package cmd

import (
	"github.com/spf13/cobra"

	"github.com/code-game-project/cli-utils/cli"

	"github.com/code-game-project/codegame-cli/games"
)

// gamesRunCmd represents the grun command
var gamesRunCmd = &cobra.Command{
	Use:                "run",
	Short:              "Download and run a game server locally from a git repository",
	DisableFlagParsing: true,
	Run: func(_ *cobra.Command, args []string) {
		installed, err := games.ListInstalled()
		checkErr("Failed to load installed games: %s", err)
		options := make([]string, 0, len(installed)+1)
		options = append(options, "Add new...")
		options = append(options, installed...)
		index := cli.Select("Select a game:", options)

		var repoURL string
		if index == 0 {
			repoURL = cli.Input("Repository URL:", true, "")
		} else {
			repoURL = options[index]
		}

		err = games.Run(repoURL, args)
		checkErr("Failed to run game server: %s", err)
	},
}

func init() {
	gamesCmd.AddCommand(gamesRunCmd)
}
