package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/code-game-project/go-utils/exec"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/semver"
	"github.com/code-game-project/go-utils/server"
	"github.com/spf13/cobra"
)

var cgDebugPath = filepath.Join(xdg.DataHome, "codegame", "bin", "cg-debug")

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "View debug logs of a game server.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		var err error
		if len(args) == 0 {
			url, err = cli.Input("Game server URL:")
			if err != nil {
				return
			}
		} else {
			url = args[0]
		}

		api, err := server.NewAPI(url)
		if err != nil {
			abort(fmt.Errorf("%s is not a CodeGame game server.", external.TrimURL(url)))
		}

		info, err := api.FetchGameInfo()
		abortf("Failed to fetch game info: %s", err)

		version, err := findDebugVersion(info.CGVersion)
		abortf("Failed to determine the correct cg-debug version to use: %s", err)

		exeName, err := installDebug(version)
		_, err = exec.Execute(false, filepath.Join(cgDebugPath, exeName), url)
		if err != nil {
			os.Exit(1)
		}
	},
}

func findDebugVersion(cgVersion string) (string, error) {
	if cgVersion == "latest" {
		version, err := external.LatestGithubTag("code-game-project", "cg-debug")
		return strings.TrimPrefix(version, "v"), err
	}

	res, err := external.LoadVersionsJSON("code-game-project", "cg-debug")
	if err != nil {
		cli.Warn("Couldn't fetch versions.json. Using latest cg-debug version.")
		version, err := external.LatestGithubTag("code-game-project", "cg-debug")
		return strings.TrimPrefix(version, "v"), err
	}

	var versions map[string]string

	err = json.Unmarshal(res, &versions)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest cg-debug version.")
		version, err := external.LatestGithubTag("code-game-project", "cg-debug")
		return strings.TrimPrefix(version, "v"), err
	}

	v := semver.CompatibleVersion(versions, cgVersion)

	if v == "latest" {
		version, err := external.LatestGithubTag("code-game-project", "cg-debug")
		return strings.TrimPrefix(version, "v"), err
	}

	v, err = external.GithubTagFromVersion("code-game-project", "cg-debug", v)
	return strings.TrimPrefix(v, "v"), err
}

func installDebug(version string) (string, error) {
	return external.InstallProgram("cg-debug", "cg-debug", "https://github.com/code-game-project/cg-debug", version, cgDebugPath)
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
