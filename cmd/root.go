package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/code-game-project/go-utils/external"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codegame",
	Short: "The official CodeGame CLI.",
	Long:  `The CodeGame CLI helps you develop CodeGame applications.`,
}

func Execute(version string) {
	rootCmd.SetVersionTemplate(`codegame-cli {{.Version}}
`)
	rootCmd.Version = version
	rootCmd.InitDefaultVersionFlag()

	if len(os.Args) <= 1 || (os.Args[1] != "upgrade" && os.Args[1] != "uninstall") {
		versionCheck(true, false)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// versionCheck returns true if a new version is available
func versionCheck(printWarning, ignoreCache bool) bool {
	if rootCmd.Version == "dev" {
		return false
	}

	latest, err := getLatestVersion(ignoreCache)
	if err != nil {
		cli.Error("Failed to fetch latest version number: %s", err)
		os.Exit(1)
	}

	if rootCmd.Version != latest {
		if printWarning {
			cli.Warn("A new version of codegame-cli is available. Run 'codegame upgrade' to install the latest version.")
		}
		return true
	}
	return false
}

func getLatestVersion(ignoreCache bool) (string, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "codegame", "cli")
	os.MkdirAll(cacheDir, 0o755)

	if ignoreCache {
		content, err := os.ReadFile(filepath.Join(cacheDir, "latest_version"))
		if err == nil {
			parts := strings.Split(string(content), "\n")
			if len(parts) >= 2 {
				cacheTime, err := strconv.Atoi(parts[0])
				if err == nil && time.Now().Unix()-int64(cacheTime) <= 60*60*3 {
					return parts[1], nil
				}
			}
		}
	}

	tag, err := external.LatestGithubTag("code-game-project", "codegame-cli")
	if err != nil {
		return "", err
	}
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), tag)), 0o644)
	return tag, nil
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

// abort prints the error to the console and terminates the program.
// abort does nothing if err is nil.
func abort(err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*exec.ExitError); !ok && !errors.Is(err, cli.ErrCanceled) {
		cli.Error(err.Error())
	}
	os.Exit(1)
}

// abortf prints the error to the console and terminates the program.
// abortf does nothing if err is nil.
func abortf(format string, err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*exec.ExitError); !ok && !errors.Is(err, cli.ErrCanceled) {
		cli.Error(format, err)
	}
	os.Exit(1)
}
