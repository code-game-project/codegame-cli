package cggenevents

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/util/exec"
	"github.com/code-game-project/codegame-cli/util/external"
)

var cgGenEventsPath = filepath.Join(xdg.DataHome, "codegame", "bin", "cg-gen-events")

// LatestCGEVersion returns the latest CGE version in the format 'x.y'.
func LatestCGEVersion() (string, error) {
	tag, err := external.LatestGithubTag("code-game-project", "cg-gen-events")
	if err != nil {
		return "", cli.Error("Couldn't determine the latest CGE version: %s", err)
	}

	return strings.TrimPrefix(strings.Join(strings.Split(tag, ".")[:2], "."), "v"), nil
}

// CGGenEvents downloads and executes the correct version of cg-gen-events.
func CGGenEvents(cgeVersion, outputDir, url, language string) error {
	exeName, err := installCGGenEvents(cgeVersion)
	if err != nil {
		return err
	}
	_, err = exec.Execute(true, filepath.Join(cgGenEventsPath, exeName), url, "--languages", language, "--output", outputDir)
	return err
}

// installCGGenEvents installs the correct version of cg-gen-events if neccessary.
func installCGGenEvents(cgeVersion string) (string, error) {
	version, err := external.GithubTagFromVersion("code-game-project", "cg-gen-events", cgeVersion)
	if err != nil {
		return "", err
	}
	version = strings.TrimPrefix(version, "v")
	return external.InstallProgram("cg-gen-events", "cg-gen-events", fmt.Sprintf("https://github.com/code-game-project/cg-gen-events"), version, cgGenEventsPath)
}

// GetEventNames uses CGGenEvents() to get a list of all the available events of the game server at url.
// It only works for CGE versions >= 0.3.
func GetEventNames(url, cgeVersion string) ([]string, error) {
	output := os.TempDir()
	err := CGGenEvents(cgeVersion, output, url, "json")
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

// GetCGEVersion returns the CGE version of the game server in the format 'x.y'.
func GetCGEVersion(baseURL string) (string, error) {
	res, err := http.Get(baseURL + "/events")
	if err != nil || res.StatusCode != http.StatusOK || (!external.HasContentType(res.Header, "text/plain") && !external.HasContentType(res.Header, "application/octet-stream")) {
		return "", cli.Error("Couldn't access /events endpoint.")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", cli.Error("Couldn't read /events file.")
	}
	return parseCGEVersion([]rune(string(data))), nil
}

func parseCGEVersion(runes []rune) string {
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
