package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-game-project/codegame-cli/external"
	"github.com/code-game-project/codegame-cli/input"
	"github.com/ogier/pflag"
)

func New() error {
	var project string
	if pflag.NArg() >= 2 {
		project = strings.ToLower(pflag.Arg(1))
	} else {
		var err error
		project, err = input.Select("Which type of project would you like to create?", []string{"Game Server", "Game Client"}, []string{"server", "client"})
		if err != nil {
			return err
		}
	}

	projectName, err := input.Input("Project name:")
	if err != nil {
		return err
	}

	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("Project '%s' already exists.", projectName)
	}

	err = os.MkdirAll(projectName, 0755)
	if err != nil {
		return err
	}

	switch project {
	case "server":
		return newServer(projectName)
	case "client":
		return newClient(projectName)
	default:
		return fmt.Errorf("Unknown project type: %s", project)
	}
}

func newServer(projectName string) error {
	return errors.New("Not implemented.")
}

func newClient(projectName string) error {
	url, err := input.Input("Enter the URL of the game server:")
	if err != nil {
		return err
	}
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "ws://") {
		url = strings.TrimPrefix(url, "ws://")
	} else if strings.HasPrefix(url, "wss://") {
		url = strings.TrimPrefix(url, "wss://")
	}
	url = strings.TrimSuffix(url, "/")
	ssl := isSSL(url)
	name, cgVersion, err := getCodeGameInfo(baseURL(url, ssl))
	if err != nil {
		return err
	}
	cgeVersion, err := getCGEVersion(baseURL(url, ssl))
	if err != nil {
		return err
	}

	var language string
	if pflag.NArg() >= 3 {
		language = strings.ToLower(pflag.Arg(2))
	} else {
		var err error
		language, err = input.Select("In which language do you want to write your project?", []string{"Go"}, []string{"go"})
		if err != nil {
			return err
		}
	}

	eventsOutput := projectName
	if language == "go" {
		eventsOutput = filepath.Join(projectName, strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", ""))
	}

	err = external.CGGenEvents(eventsOutput, baseURL(url, ssl), cgeVersion, language)
	if err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31mFailed to generate event definitions:", err, "\x1b[0m")
	}

	switch language {
	case "go":
		err = newClientGo(projectName, url, cgVersion)
	default:
		return fmt.Errorf("Unsupported language: %s", language)
	}
	if err != nil {
		return err
	}

	fmt.Printf("\x1b[32mSuccessfully created project in '%s'.\n\x1b[0m", projectName)

	return nil
}

func getCodeGameInfo(baseURL string) (string, string, error) {
	type response struct {
		Name      string `json:"name"`
		CGVersion string `json:"cg_version"`
	}
	res, err := http.Get(baseURL + "/info")
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", "", err
	}

	return data.Name, data.CGVersion, nil
}

func getCGEVersion(baseURL string) (string, error) {
	res, err := http.Get(baseURL + "/events")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
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

func baseURL(domain string, ssl bool) string {
	if ssl {
		return "https://" + domain
	} else {
		return "http://" + domain
	}
}

func isSSL(domain string) bool {
	res, err := http.Get("https://" + domain)
	if err == nil {
		res.Body.Close()
		return true
	}
	return false
}
