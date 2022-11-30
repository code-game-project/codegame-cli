package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/server"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// gameCreateCmd represents the game create command
var gameCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new game on the a server.",
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

		type request struct {
			Public    bool `json:"public"`
			Protected bool `json:"protected"`
			Config    any  `json:"config,omitempty"`
		}
		public, err := cmd.Flags().GetBool("public")
		abort(err)
		protected, err := cmd.Flags().GetBool("protected")
		abort(err)

		data, err := json.Marshal(request{
			Public:    public,
			Protected: protected,
		})
		abort(err)

		body := bytes.NewBuffer(data)
		resp, err := http.Post(api.BaseURL()+"/games", "application/json", body)
		abort(err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			abort(fmt.Errorf("invalid response code: expected: %d, got: %d", http.StatusCreated, resp.StatusCode))
		}

		type response struct {
			GameId     string `json:"game_id"`
			JoinSecret string `json:"join_secret"`
		}
		var r response
		err = json.NewDecoder(resp.Body).Decode(&r)
		abort(err)

		out := colorable.NewColorableStdout()
		fmt.Fprintf(out, "%sGame ID:%s %s\n", cli.Cyan, cli.Reset, r.GameId)
		if r.JoinSecret != "" {
			fmt.Fprintf(out, "%sJoin secret:%s %s\n", cli.Cyan, cli.Reset, r.JoinSecret)
		}
	},
}

func init() {
	gameCmd.AddCommand(gameCreateCmd)
	gameCreateCmd.Flags().BoolP("public", "", false, "The game is displayed on a public game list.")
	gameCreateCmd.Flags().BoolP("protected", "", false, "You can only join the game with the returned join secret.")
}
