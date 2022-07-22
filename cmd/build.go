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

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "Build the current project.",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		rootRelative, err := cgfile.FindProjectRootRelative()
		abort(err)

		data, err := cgfile.LoadCodeGameFile(rootRelative)
		abortf("failed to load .codegame.json: %w", err)

		cmdArgs := []string{"build"}
		cmdArgs = append(cmdArgs, args...)

		switch data.Lang {
		case "go":
			err = modules.Execute("go", "latest", data.Type, cmdArgs...)
			abort(err)
		default:
			abort(fmt.Errorf("'build' is not supported for '%s'", data.Lang))
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
