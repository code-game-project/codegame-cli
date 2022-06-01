package external

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/cli"
)

const currentVersion = "0.4.1"

func CheckVersion() {
	latest, err := getLatestVersion()
	if err != nil {
		return
	}
	if currentVersion != latest {
		cli.Warn("You are using an old version of codegame-cli (v%s).\nUpdate to the latest version (v%s): https://github.com/code-game-project/codegame-cli#installation", currentVersion, latest)
	}
}

func getLatestVersion() (string, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "codegame", "cli")
	os.MkdirAll(cacheDir, 0755)

	content, err := os.ReadFile(filepath.Join(cacheDir, "latest_version"))
	if err == nil {
		parts := strings.Split(string(content), "\n")
		if len(parts) >= 2 {
			cacheTime, err := strconv.Atoi(parts[0])
			if err == nil && time.Now().Unix()-int64(cacheTime) <= 3*24*60*60 {
				return parts[1], nil
			}
		}
	}

	tag, err := LatestGithubTag("code-game-project", "codegame-cli")
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(tag, "v")
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), version)), 0644)
	return version, nil
}
