package cmd

import (
	"fmt"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/server"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// gameListCmd represents the game list command
var gameListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all public games of a game server.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var gameURL string
		var err error
		if len(args) > 0 {
			gameURL = args[0]
		} else if gameURL = findGameURL(); gameURL == "" {
			gameURL, err = cli.Input("Game URL:")
			abort(err)
		}

		api, err := server.NewAPI(gameURL)
		abort(err)

		protected, err := cmd.Flags().GetBool("protected")
		abort(err)
		unprotected, err := cmd.Flags().GetBool("unprotected")
		abort(err)
		private, public, err := api.ListGames(unprotected, protected)
		abort(err)

		out := colorable.NewColorableStdout()
		fmt.Fprintf(out, "%sPrivate:%s %d\n", cli.Cyan, cli.Reset, private)
		if len(public) == 0 {
			fmt.Fprintf(out, "%sPublic:%s none\n", cli.Cyan, cli.Reset)
		} else {
			cli.PrintColor(cli.Cyan, "Public:")
			for _, g := range public {
				if g.Protected {
					cli.Print("- %s (%d players, protected)", g.Id, g.Players)
				}
				if !g.Protected {
					cli.Print("- %s (%d players)", g.Id, g.Players)
				}
			}
		}
	},
}

func init() {
	gameCmd.AddCommand(gameListCmd)
	gameListCmd.Flags().BoolP("protected", "", false, "Only show protected games.")
	gameListCmd.Flags().BoolP("unprotected", "", false, "Only show unprotected games.")
}
