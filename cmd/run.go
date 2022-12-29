package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/config"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/modules"
	"github.com/code-game-project/go-utils/semver"
	"github.com/code-game-project/go-utils/server"
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

		if data.GameVersion != "" {
			wrapMaj, wrapMin, _, err := semver.ParseVersion(data.GameVersion)
			if err != nil {
				goto skipGameVersionCheck
			}

			api, err := server.NewAPI(data.URL)
			if err != nil {
				goto skipGameVersionCheck
			}
			info, err := api.FetchGameInfo()
			if err != nil || info.Version == "" {
				goto skipGameVersionCheck
			}

			gameMaj, gameMin, _, err := semver.ParseVersion(info.Version)
			if err != nil {
				goto skipGameVersionCheck
			}

			if wrapMaj != gameMaj || wrapMin != gameMin {
				cli.Warn("Game version mismatch. Server: v%s, client: v%s. Please run 'codegame update'.", info.Version, data.GameVersion)
			}
		}
	skipGameVersionCheck:

		runData := modules.RunData{
			Lang: data.Lang,
			Args: args,
		}

		if _, ok := os.LookupEnv("CG_PORT"); !ok {
			conf := config.Load()
			port := findAvailablePort(conf.DevPort)
			os.Setenv("CG_PORT", fmt.Sprintf("%d", port))
		}

		switch data.Lang {
		case "cs", "go", "java", "js", "ts":
			err = modules.ExecuteRun(runData, data)
			abort(err)
		default:
			abort(fmt.Errorf("'run' is not supported for '%s'", data.Lang))
		}
	},
}

func findAvailablePort(port int) int {
	for i := port; i < port+100; i++ {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			listener.Close()
			return i
		}
	}
	return port
}

func init() {
	rootCmd.AddCommand(runCmd)
}
