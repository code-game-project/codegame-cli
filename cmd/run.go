package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/code-game-project/cli-utils/cgfile"

	"github.com/code-game-project/codegame-cli/run"
)

var (
	runCmdSpectate bool
	runCmdPort     int
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Execute a CodeGame application",
	DisableFlagParsing: true,
	Run: func(_ *cobra.Command, args []string) {
		cgFile, err := cgfile.Load("")
		checkErr("Not in a CodeGame project directory: %s", err)
		switch cgFile.ProjectType {
		case "client":
			err = run.RunClient(cgFile, runCmdSpectate, args)
		case "server":
			err = run.RunServer(cgFile, runCmdPort, args)
		default:
			err = errors.New("unknown project type")
		}
		checkErr("Failed to execute project: %s", err)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVar(&runCmdSpectate, "spectate", false, "Spectate the game (only available for clients)")
	runCmd.Flags().IntVar(&runCmdPort, "port", 0, "The port to listen on (only available for servers)")
}
