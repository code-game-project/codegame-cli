package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Bananenpro/cli"
	"github.com/Bananenpro/pflag"
	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/commands"
	"github.com/code-game-project/codegame-cli/util/external"
	"github.com/code-game-project/codegame-cli/util/semver"
)

const Version = "0.6.0"

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nDescription:")
		fmt.Fprintln(os.Stderr, "The official CodeGame CLI.")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "\tnew \tCreate a new project.")
		fmt.Fprintln(os.Stderr, "\trun \tRun a project.")
		fmt.Fprintln(os.Stderr, "\tinfo \tDisplay some info about a game server.")
		fmt.Fprintln(os.Stderr, "\tdocs \tOpen the CodeGame documentation in a web browser.")
		fmt.Fprintln(os.Stderr, "\nAbout: https://code-game.org")
		fmt.Fprintln(os.Stderr, "Copyright (c) 2022 CodeGame Contributors (https://code-game.org/contributors)")
		pflag.PrintDefaults()
	}

	if len(os.Args) == 1 {
		pflag.Usage()
		os.Exit(1)
	}

	checkVersion()

	command := strings.ToLower(os.Args[1])

	var err error
	switch command {
	case "new":
		err = commands.New()
	case "info":
		err = commands.Info()
	case "docs":
		err = commands.Docs()
	case "run":
		err = commands.Run()
	default:
		cli.Error("Unknown command: %s", strings.ToLower(pflag.Arg(0)))
		pflag.Usage()
		os.Exit(1)
	}
	if err != nil {
		cli.CancelLoading()
		cli.CancelProgressBar()
		os.Exit(1)
	}
}

// checkVersion prints a warning, if there is a newer version of codegame-cli available.
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
		cli.Warn("You are using an old version of codegame-cli (v%s).\nUpdate to the latest version (v%s): https://github.com/code-game-project/codegame-cli#installation", Version, latest)
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

	tag, err := external.LatestGithubTag("code-game-project", "codegame-cli")
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(tag, "v")
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), version)), 0644)
	return version, nil
}
