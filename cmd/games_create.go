package cmd

import (
	"encoding/json"

	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/request"
	"github.com/code-game-project/cli-utils/server"
	"github.com/spf13/cobra"
)

// gamesCreateCmd represents the games create command
var gamesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a game on a game server",
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

		public, err := cmd.Flags().GetBool("public")
		checkErr("%s", err)
		protected, err := cmd.Flags().GetBool("protected")
		checkErr("%s", err)
		configStr, err := cmd.Flags().GetString("config")
		checkErr("%s", err)

		var config json.RawMessage
		err = json.Unmarshal([]byte(configStr), &config)
		checkErr("Invalid config: %s", err)

		gameID, joinSecret, err := server.CreateGame(url, public, protected, config)
		checkErr("Failed to create game: %s", err)
		if joinSecret != "" {
			cli.Print("%sGame ID:%s     %s", cli.Cyan, cli.Reset, gameID)
			cli.Print("%sJoin secret:%s %s", cli.Cyan, cli.Reset, joinSecret)
		} else {
			cli.Print("%sGame ID:%s %s", cli.Cyan, cli.Reset, gameID)
		}
	},
}

func init() {
	gamesCmd.AddCommand(gamesCreateCmd)
	gamesCreateCmd.Flags().Bool("public", false, "Create a public game")
	gamesCreateCmd.Flags().Bool("protected", false, "Create a protected game")
	gamesCreateCmd.Flags().String("config", "null", "The game configuration (json)")
}
