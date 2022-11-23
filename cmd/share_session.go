package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/config"
	"github.com/code-game-project/go-utils/external"
	"github.com/spf13/cobra"
)

// shareSessionCmd represents the share session command
var shareSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Share a session with CodeGame Share (same as 'codegame session export').",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		session, err := selectSession(args)
		abortf("Failed to load session: %s", err)

		password, err := cli.InputOptional("Password (optional):")
		abort(err)

		type reqSession struct {
			GameId       string `json:"game_id"`
			PlayerId     string `json:"player_id"`
			PlayerSecret string `json:"player_secret"`
		}

		type request struct {
			GameURL  string     `json:"game_url"`
			Username string     `json:"username"`
			Session  reqSession `json:"session"`
			Password string     `json:"password"`
		}

		data := request{
			GameURL:  session.GameURL,
			Username: session.Username,
			Session: reqSession{
				GameId:       session.GameId,
				PlayerId:     session.PlayerId,
				PlayerSecret: session.PlayerSecret,
			},
			Password: password,
		}

		jsonData, err := json.Marshal(data)
		abort(err)

		conf := config.Load()
		shareURL := external.TrimURL(conf.ShareURL)
		baseURL := external.BaseURL("http", external.IsTLS(shareURL), shareURL)

		resp, err := http.Post(baseURL+"/session", "application/json", bytes.NewBuffer(jsonData))
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
			if err != nil {
				cli.Error("Failed to decode error message: %s", err)
			} else {
				cli.Error(res.Error)
			}
			return
		}

		type response struct {
			Id string `json:"id"`
		}
		var res response
		err = json.NewDecoder(resp.Body).Decode(&res)
		abortf("Failed to decode server response: %s", err)
		cli.Success("Success! You can import your session on another device with the following command:")
		cli.PrintColor(cli.Cyan, "codegame session import %s", res.Id)
	},
}

func init() {
	shareCmd.AddCommand(shareSessionCmd)
}
