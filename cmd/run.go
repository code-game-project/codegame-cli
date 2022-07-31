package cmd

import (
	"fmt"
	"os"

	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/modules"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run the current project.",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := cgfile.FindProjectRoot()
		abort(err)
		err = os.Chdir(root)
		abort(err)

		data, err := cgfile.LoadCodeGameFile("")
		abortf("failed to load .codegame.json: %w", err)
		data.URL = external.TrimURL(data.URL)
		abort(data.Write(""))

		runData := modules.RunData{
			Lang: data.Lang,
			Args: args,
		}

		switch data.Lang {
		case "go":
			err = modules.ExecuteRun(runData, data)
			abort(err)
		default:
			abort(fmt.Errorf("'run' is not supported for '%s'", data.Lang))
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
