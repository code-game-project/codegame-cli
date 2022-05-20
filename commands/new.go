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

	switch project {
	case "server":
		return newServer()
	case "client":
		return newClient()
	default:
		return fmt.Errorf("Unknown project type: %s", project)
	}
}

func newServer() error {
	return errors.New("Not implemented.")
}

func newClient() error {
	domain, err := input.Input("Enter the domain of the game server:")
	if err != nil {
		return err
	}
	if strings.HasPrefix(domain, "http://") {
		domain = strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(domain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
	} else if strings.HasPrefix(domain, "ws://") {
		domain = strings.TrimPrefix(domain, "ws://")
	} else if strings.HasPrefix(domain, "wss://") {
		domain = strings.TrimPrefix(domain, "wss://")
	}
	domain = strings.TrimSuffix(domain, "/")
	ssl := isSSL(domain)
	name, _, err := getCodeGameInfo(baseURL(domain, ssl))
	if err != nil {
		return err
	}
	cgeVersion, err := getCGEVersion(baseURL(domain, ssl))
	if err != nil {
		return err
	}

	projectName, err := input.Input("Project name:")
	if err != nil {
		return err
	}

	err = os.MkdirAll(projectName, 0755)
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
		eventsOutput = filepath.Join(projectName, name)
	}

	err = external.CGGenEvents(eventsOutput, baseURL(domain, ssl), cgeVersion, language)
	if err != nil {
		return err
	}
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
