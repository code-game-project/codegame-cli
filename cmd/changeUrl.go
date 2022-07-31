package cmd

import (
	"errors"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/server"
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
		if len(args) > 0 {
			url = strings.ToLower(args[0])
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

		config.URL = external.TrimURL(url)
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
