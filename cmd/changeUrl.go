/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/pkg/cgfile"
	"github.com/code-game-project/codegame-cli/pkg/server"
	"github.com/spf13/cobra"
)

// changeUrlCmd represents the changeUrl command
var changeUrlCmd = &cobra.Command{
	Use:   "changeUrl",
	Short: "Permanently switch to a different game URL.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		root, err := cgfile.FindProjectRoot()
		abort(err)

		var url string
		if len(os.Args) >= 3 {
			url = strings.ToLower(os.Args[2])
		} else {
			var err error
			url, err = cli.Input("New game URL:")
			abort(err)
		}
		api, err := server.NewAPI(url)
		abort(err)

		config, err := cgfile.LoadCodeGameFile(root)
		abort(err)

		if config.Type != "client" {
			abort(errors.New("project is not a client"))
		}

		info, err := api.FetchGameInfo()
		abort(err)
		if info.Name != config.Game {
			abort(errors.New("The URL points to a different game."))
		}

		prevURL := config.URL

		config.URL = url
		err = config.Write(root)
		abort(err)

		err = update()
		if err != nil {
			config.URL = prevURL
			config.Write(root)
			abort(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(changeUrlCmd)
}
