package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

		output, err := cmd.Flags().GetString("output")
		abort(err)

		os, err := cmd.Flags().GetString("os")
		abort(err)
		os = strings.ToLower(os)
		if os != "current" && os != "windows" && os != "macos" && os != "linux" {
			abort(fmt.Errorf("OS '%s' is not supported. (possible values: windows, macos, linux)", os))
		}
		arch, err := cmd.Flags().GetString("arch")
		abort(err)
		arch = strings.ToLower(arch)
		if arch != "current" && arch != "x64" && arch != "x86" && arch != "arm32" && arch != "arm64" {
			abort(fmt.Errorf("Architecture '%s' is not supported. (possible values: x64, x86, arm32, arm64)", arch))
		}

		if os == "macos" && arch == "arm32" {
			abort(errors.New("macOS does not support arm32. Try arm64 instead."))
		}

		buildData := modules.BuildData{
			Lang:   data.Lang,
			Output: output,
			OS:     os,
			Arch:   arch,
		}
		switch data.Lang {
		case "cs", "go", "js", "ts":
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
	buildCmd.Flags().StringP("os", "", "current", "The target OS for compiled languages. (possible values: windows, macos, linux)")
	buildCmd.Flags().StringP("arch", "", "current", "The target architecture for compiled languages. (possible values: x64, x86, arm32, arm64)")
}
