/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/code-game-project/codegame-cli/pkg/cgfile"
	"github.com/code-game-project/codegame-cli/pkg/modules"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run the current project.",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootRelative, err := cgfile.FindProjectRootRelative()
		abort(err)

		data, err := cgfile.LoadCodeGameFile(rootRelative)
		abortf("failed to load .codegame.json: %w", err)

		cmdArgs := []string{"run"}
		cmdArgs = append(cmdArgs, args...)

		switch data.Lang {
		case "go":
			err = modules.Execute("go", "latest", data.Type, cmdArgs...)
			abort(err)
		default:
			abort(fmt.Errorf("'run' is not supported for '%s'", data.Lang))
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
