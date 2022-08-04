package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// shareGameCmd represents the share game command
var shareGameCmd = &cobra.Command{
	Use:   "game",
	Short: "Share a game with share.code-game.org.",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var gameURL string
		var gameId string

		if len(args) > 0 {
			gameURL = args[0]
			if len(args) > 1 {
				gameId = args[1]
			}
		}

		var err error
		if gameURL == "" {
			fromSession, err := cli.YesNo("Select game URL from session?", true)
			abort(err)
			if fromSession {
				urls, err := sessions.ListGames()
				abortf("Failed to load games: %s", err)
				selected, err := cli.Select("Game URL:", urls)
				abort(err)
				gameURL = urls[selected]
			} else {
				gameURL, err = cli.Input("Game URL:")
				abort(err)
			}
		}

		if gameId == "" {
			fromSession, err := cli.YesNo("Select game ID from session?", true)
			abort(err)
			if fromSession {
				usernames, err := sessions.ListUsernames(gameURL)
				abortf("Failed to load usernames: %s", err)
				selected, err := cli.Select("Username:", usernames)
				abort(err)
				session, err := sessions.LoadSession(gameURL, usernames[selected])
				abortf("Failed to load session: %s", err)
				gameId = session.GameId
			} else {
				gameId, err = cli.Input("Game ID:")
				abort(err)
			}
		}

		type request struct {
			GameURL string `json:"game_url"`
			GameId  string `json:"game_id"`
		}

		data := request{
			GameURL: gameURL,
			GameId:  gameId,
		}

		jsonData, err := json.Marshal(data)
		abort(err)

		resp, err := http.Post("https://share.code-game.org/game", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			cli.Error("Failed to upload data: %s", err)
			return
		}
		if resp.StatusCode != http.StatusCreated {
			type response struct {
				Error string `json:"error"`
			}
			var res response
			err = json.NewDecoder(resp.Body).Decode(&res)
			cli.Error(res.Error)
			return
		}

		type response struct {
			Id string `json:"id"`
		}
		var res response
		err = json.NewDecoder(resp.Body).Decode(&res)
		abortf("Failed to decode server response: %s", err)
		cli.Success("Success! You can view the game details with the following link:")
		cli.PrintColor(cli.Cyan, "https://share.code-game.org/%s", res.Id)
	},
}

func init() {
	shareCmd.AddCommand(shareGameCmd)
}
