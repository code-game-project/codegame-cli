package external

import (
	"encoding/json"
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
		return "", cli.Error("Couldn't access git tags from 'github.com/%s/%s'.", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", cli.Error("Couldn't decode git tag data.")
	}
	return data[0].Name, nil
}

func LatestCGEVersion() (string, error) {
	tag, err := LatestGithubTag("code-game-project", "cg-gen-events")
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(strings.Join(strings.Split(tag, ".")[:2], "."), "v"), nil
}

func InstallCGGenEvents(cgeVersion string) error {
	exeName := fmt.Sprintf("cg-gen-events-%s", strings.ReplaceAll(cgeVersion, ".", "_"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cgGenEventsPath, exeName)); err == nil {
		return nil
	}

	version, err := GithubTagFromVersion("code-game-project", "cg-gen-events", cgeVersion)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("cg-gen-events-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		filename = fmt.Sprintf("cg-gen-events-%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	}

	res, err := http.Get(fmt.Sprintf("https://github.com/code-game-project/cg-gen-events/releases/download/%s/%s", version, filename))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = os.MkdirAll(cgGenEventsPath, 0755)
	if err != nil {
		return err
	}

	return UntargzFile(res.Body, "cg-gen-events", filepath.Join(cgGenEventsPath, exeName))
}

func CGGenEvents(outputDir, url, cgeVersion, language string) error {
	exeName := fmt.Sprintf("cg-gen-events-%s", strings.ReplaceAll(cgeVersion, ".", "_"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cgGenEventsPath, exeName)); err != nil {
		err = InstallCGGenEvents(cgeVersion)
		if err != nil {
			return err
		}
	}

	return Execute(filepath.Join(cgGenEventsPath, exeName), url, "--languages", language, "--output", outputDir)
}
