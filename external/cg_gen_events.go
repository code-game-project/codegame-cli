package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// requires cgeVersion >= 0.3
func GetEventNames(url, cgeVersion string) ([]string, error) {
	output := os.TempDir()
	err := CGGenEvents(output, url, cgeVersion, "json")
	if err != nil {
		return nil, err
	}

	type event struct {
		Name string `json:"name"`
	}
	type data struct {
		Events []event `json:"events"`
	}

	path := filepath.Join(output, "events.json")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)
	defer file.Close()

	var object data
	err = json.NewDecoder(file).Decode(&object)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(object.Events))
	for i, event := range object.Events {
		names[i] = event.Name
	}
	return names, nil
}

func GetCGEVersion(baseURL string) (string, error) {
	res, err := http.Get(baseURL + "/events")
	if err != nil || res.StatusCode != http.StatusOK || !HasContentType(res.Header, "text/plain") {
		return "", cli.Error("Couldn't access /events endpoint.")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", cli.Error("Couldn't read /events file.")
	}
	return ParseCGEVersion([]rune(string(data))), nil
}

func ParseCGEVersion(runes []rune) string {
	index := 0
	commentNestingLevel := 0
	for index < len(runes) && (runes[index] == ' ' || runes[index] == '\r' || runes[index] == '\n' || runes[index] == '\t' || (index < len(runes)-1 && runes[index] == '/' && runes[index+1] == '*') || (index < len(runes)-1 && runes[index] == '*' && runes[index+1] == '/') || (index < len(runes)-1 && runes[index] == '/' && runes[index+1] == '/') || commentNestingLevel > 0) {
		if runes[index] == '/' {
			if runes[index+1] == '/' {
				for index < len(runes) && runes[index] != '\n' {
					index++
				}
			} else {
				commentNestingLevel++
			}
		}
		if runes[index] == '*' {
			commentNestingLevel--
		}
		index++
	}

	words := strings.Fields(string(runes[index:]))
	for i, w := range words {
		if w == "version" && i < len(words)-1 {
			return words[i+1]
		}
	}

	return ""
}
