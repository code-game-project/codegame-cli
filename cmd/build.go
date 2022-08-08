package cmd

import (
	"fmt"
	"os"

	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/modules"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the current project.",
	Run: func(cmd *cobra.Command, args []string) {
		root, err := cgfile.FindProjectRoot()
		abort(err)
		err = os.Chdir(root)
		abort(err)

		data, err := cgfile.LoadCodeGameFile("")
		abortf("failed to load .codegame.json: %w", err)
		data.URL = external.TrimURL(data.URL)
		abort(data.Write(""))

		cmdArgs := []string{"build"}
		cmdArgs = append(cmdArgs, args...)

		output, err := cmd.Flags().GetString("output")
		abort(err)
		buildData := modules.BuildData{
			Lang:   data.Lang,
			Output: output,
		}
		switch data.Lang {
		case "go", "js", "ts":
			err = modules.ExecuteBuild(buildData, data)
			abort(err)
		default:
			abort(fmt.Errorf("'build' is not supported for '%s'", data.Lang))
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("output", "o", "", "The name of the output file.")
}
