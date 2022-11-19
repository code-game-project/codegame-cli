package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/config"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/sessions"
	"github.com/spf13/cobra"
)

// sessionImportCmd represents the session import command
var sessionImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a session from CodeGame Share.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		var err error
		if len(args) > 0 {
			id = args[0]
		} else {
			id, err = cli.Input("ID:")
			if err != nil {
				return
			}
		}

		conf := config.Load()
		shareURL := external.TrimURL(conf.ShareURL)
		baseURL := external.BaseURL("http", external.IsTLS(shareURL), shareURL)

		resp, err := http.Get(fmt.Sprintf(baseURL+"/%s?type=session", id))
		abortf(fmt.Sprintf("Failed to contact %s: %s", conf.ShareURL, "%s"), err)

		type resSession struct {
			GameId       string `json:"game_id"`
			PlayerId     string `json:"player_id"`
			PlayerSecret string `json:"player_secret"`
		}

		type response struct {
			Error    string     `json:"error"`
			GameURL  string     `json:"game_url"`
			Username string     `json:"username"`
			Session  resSession `json:"session"`
		}

		var data response
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			abortf("Failed to decode server response: %s", err)
		}
		if data.Error != "" {
			abort(fmt.Errorf(data.Error))
		}

		session := sessions.NewSession(data.GameURL, data.Username, data.Session.GameId, data.Session.PlayerId, data.Session.PlayerSecret)
		err = session.Save()
		abortf("Failed to save session: %s", err)

		cli.Success("Successfully imported %s@%s!", session.Username, session.GameURL)
	},
}

func init() {
	sessionCmd.AddCommand(sessionImportCmd)
}
