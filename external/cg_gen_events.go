package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/cli"
)

var cgGenEventsPath = filepath.Join(xdg.DataHome, "codegame", "bin", "cg-gen-events")

func LatestGithubTag(owner, repo string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK || !HasContentType(res.Header, "application/json") {
		return "", fmt.Errorf("failed to access git tags from 'github.com/%s/%s'.", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", errors.New("failed to decode git tag data.")
	}
	return data[0].Name, nil
}

func LatestCGEVersion() (string, error) {
	tag, err := LatestGithubTag("code-game-project", "cg-gen-events")
	if err != nil {
		return "", cli.Error("Couldn't determine the latest CGE version: %s", err)
	}

	return strings.TrimPrefix(strings.Join(strings.Split(tag, ".")[:2], "."), "v"), nil
}

func InstallCGGenEvents(version string) (string, error) {
	exeName := fmt.Sprintf("cg-gen-events_%s", strings.ReplaceAll(version, ".", "-"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cgGenEventsPath, exeName)); err == nil {
		return exeName, nil
	}

	filename := fmt.Sprintf("cg-gen-events-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		filename = fmt.Sprintf("cg-gen-events-%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	}

	res, err := http.Get(fmt.Sprintf("https://github.com/code-game-project/cg-gen-events/releases/download/v%s/%s", version, filename))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = os.MkdirAll(cgGenEventsPath, 0755)
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		return exeName, UnzipFile(res.Body, "cg-gen-events.exe", filepath.Join(cgGenEventsPath, exeName))
	}
	return exeName, UntargzFile(res.Body, "cg-gen-events", filepath.Join(cgGenEventsPath, exeName))
}

func CGGenEvents(outputDir, url, cgeVersion, language string) error {
	version, err := GithubTagFromVersion("code-game-project", "cg-gen-events", cgeVersion)
	if err != nil {
		return err
	}
	version = strings.TrimPrefix(version, "v")

	exeName, err := InstallCGGenEvents(version)
	if err != nil {
		return err
	}

	binaries, err := os.ReadDir(cgGenEventsPath)
	if err != nil {
		return err
	}
	for _, b := range binaries {
		info, err := b.Info()
		if err == nil && info.Name() != exeName && strings.HasPrefix(info.Name(), fmt.Sprintf("cg-gen-events_%s", strings.ReplaceAll(cgeVersion, ".", "-"))) {
			os.Remove(filepath.Join(cgGenEventsPath, info.Name()))
		}
	}

	out, err := ExecuteHidden(filepath.Join(cgGenEventsPath, exeName), url, "--languages", language, "--output", outputDir)
	if err != nil {
		cli.Error(out)
	}
	return err
}
