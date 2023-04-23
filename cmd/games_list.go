package cmd

import (
	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/request"
	"github.com/code-game-project/cli-utils/server"
	"github.com/spf13/cobra"
)

// gamesListCmd represents the games list command
var gamesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all public games on a game server",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 0 {
			file, err := cgfile.Load("")
			if err == nil {
				url = file.GameURL
			}
			url = cli.Input("Game URL:", true, url)
		} else {
			url = args[0]
		}
		url = request.TrimURL(url)

		protected, err := cmd.Flags().GetBool("protected")
		checkErr("%s", err)

		private, public, err := server.FetchGames(url, protected)
		checkErr("Failed to fetch games: %s", err)
		cli.Print("%sPrivate:%s %d", cli.Cyan, cli.Reset, private)
		if len(public) == 0 {
			cli.Print("%sPublic:%s  0", cli.Cyan, cli.Reset)
		} else {
			cli.PrintColor(cli.Cyan, "Public:")
			for _, g := range public {
				if g.Protected {
					cli.Print("- %s (%d players, protected)", g.ID, g.Players)
				} else {
					cli.Print("- %s (%d players)", g.ID, g.Players)
				}
			}
		}
	},
}

func init() {
	gamesCmd.AddCommand(gamesListCmd)
	gamesListCmd.Flags().Bool("protected", false, "Include protected games")
}
