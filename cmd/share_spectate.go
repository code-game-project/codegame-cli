package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/spf13/cobra"
)

// shareSpectateCmd represents the share spectate command
var shareSpectateCmd = &cobra.Command{
	Use:   "spectate",
	Short: "Share a spectate link with share.code-game.org.",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var gameURL string
		var gameId string
		var playerId string
		var playerSecret string

		fromSession, err := cli.YesNo("Select from session?", true)
		abort(err)
		if fromSession {
			session, err := selectSession(args)
			abortf("Failed to load session: %s", err)
			gameURL = session.GameURL
			gameId = session.GameId
			playerId = session.PlayerId
			playerSecret = session.PlayerSecret
		} else {
			if len(args) > 0 {
				gameURL = args[0]
			} else if gameURL = findGameURL(); gameURL != "" {
				cli.Print("Game URL: %s", gameURL)
			} else {
				gameURL, err = cli.Input("Game URL:")
				abort(err)
			}
			gameId, err = cli.Input("Game ID:")
			abort(err)
			playerId, err = cli.Input("Player ID:")
			abort(err)
			playerSecret, err = cli.Input("Player Secret:")
			abort(err)
		}

		type request struct {
			GameURL      string `json:"game_url"`
			GameId       string `json:"game_id"`
			PlayerId     string `json:"player_id"`
			PlayerSecret string `json:"player_secret"`
		}

		data := request{
			GameURL:      gameURL,
			GameId:       gameId,
			PlayerId:     playerId,
			PlayerSecret: playerSecret,
		}

		jsonData, err := json.Marshal(data)
		abort(err)

		resp, err := http.Post("https://share.code-game.org/spectate", "application/json", bytes.NewBuffer(jsonData))
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
		cli.Success("Success! You can spectate the game with the following link:")
		cli.PrintColor(cli.Cyan, "https://share.code-game.org/%s", res.Id)
	},
}

func init() {
	shareCmd.AddCommand(shareSpectateCmd)
}
