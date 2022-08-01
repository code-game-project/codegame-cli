package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// sessionExportCmd represents the session export command
var sessionExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a session to share.code-game.org.",
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
		abortf("Failed to load session: %s", err)

		type reqSession struct {
			GameId       string `json:"game_id"`
			PlayerId     string `json:"player_id"`
			PlayerSecret string `json:"player_secret"`
		}

		type request struct {
			GameURL  string     `json:"game_url"`
			Username string     `json:"username"`
			Session  reqSession `json:"session"`
		}

		data := request{
			GameURL:  gameURL,
			Username: username,
			Session: reqSession{
				GameId:       session.GameId,
				PlayerId:     session.PlayerId,
				PlayerSecret: session.PlayerSecret,
			},
		}

		jsonData, err := json.Marshal(data)
		abort(err)

		resp, err := http.Post("https://share.code-game.org/session", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			cli.Error("Failed to upload session: %s", err)
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
		abortf("Failed to decode session data from server: %s", err)
		cli.Success("Success! You can import your session with the following command:")
		cli.PrintColor(cli.Cyan, "codegame session import %s", res.Id)
	},
}

func init() {
	sessionCmd.AddCommand(sessionExportCmd)
}
