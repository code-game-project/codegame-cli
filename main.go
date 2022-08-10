package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/semver"

	cgExec "github.com/code-game-project/go-utils/exec"

	"github.com/code-game-project/codegame-cli/cmd"
)

const Version = "0.7.1"

func main() {
	checkVersion()
	cmd.Execute(Version)
}

// checkVersion prints a warning, if there is a newer version of codegame-cli available.
// On macOS and linux the user is offered to update automatically.
func checkVersion() {
	latest, err := getLatestVersion()
	if err != nil {
		return
	}

	currentMajor, currentMinor, currentPatch, err := semver.ParseVersion(Version)
	if err != nil {
		return
	}

	latestMajor, latestMinor, latestPatch, err := semver.ParseVersion(latest)
	if err != nil {
		return
	}

	if latestMajor > currentMajor || (latestMajor == currentMajor && latestMinor > currentMinor) || (latestMajor == currentMajor && latestMinor == currentMinor && latestPatch > currentPatch) {
		_, shErr := exec.LookPath("sh")
		_, sudoErr := exec.LookPath("sudo")
		_, curlErr := exec.LookPath("curl")
		_, tarErr := exec.LookPath("tar")
		cgBin, codegameErr := os.Stat("/usr/local/bin/codegame")
		if codegameErr == nil && !cgBin.IsDir() && shErr == nil && sudoErr == nil && curlErr == nil && tarErr == nil && (runtime.GOOS == "darwin" || runtime.GOOS == "linux") {
			update()
		} else {
			cli.Warn("You are using an old version of codegame-cli (v%s).\nUpdate to the latest version (v%s): https://github.com/code-game-project/codegame-cli#installation", Version, latest)
		}
	}
}

func update() {
	yes, err := cli.YesNo("A new version is available. Do you want to update now?", true)
	if err != nil {
		os.Exit(0)
	}
	if !yes {
		return
	}

	_, err = cgExec.Execute(false, "sh", "-c", fmt.Sprintf("curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-%s-%s.tar.gz | tar -xz codegame && sudo mv codegame /usr/local/bin", runtime.GOOS, runtime.GOARCH))
	if err != nil {
		cli.Error("Update failed.")
		os.Exit(1)
	}
	cli.Success("Update successful.")
	os.Exit(0)
}

func getLatestVersion() (string, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "codegame", "cli")
	os.MkdirAll(cacheDir, 0o755)

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

	tag, err := external.LatestGithubTag("code-game-project", "codegame-cli")
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(tag, "v")
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), version)), 0o644)
	return version, nil
}
