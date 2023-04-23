package cmd

import (
	"os"

	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/components"
	"github.com/code-game-project/cli-utils/exec"
	"github.com/code-game-project/cli-utils/feedback"
	"github.com/code-game-project/cli-utils/request"
	"github.com/code-game-project/cli-utils/server"
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "View debug logs of a game server",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(_ *cobra.Command, args []string) {
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
		feedback.Info("codegame-cli", "Debugging %s...", url)

		info, err := server.FetchGameInfo(url)
		if err != nil {
			feedback.Fatal("codegame-cli", "%s is not a valid CodeGame game server or is not reachable", url)
			os.Exit(1)
		}
		cgDebug, err := components.CGDebug(info.CGVersion)
		checkErr("Failed execute cg-debug: %w", err)
		err = exec.Execute(cgDebug, url)
		checkErr("Failed to debug game server: %w", err)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
