package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/external"
	"github.com/ogier/pflag"
)

func New() error {
	var project string
	if pflag.NArg() >= 2 {
		project = strings.ToLower(pflag.Arg(1))
	} else {
		var err error
		project, err = cli.Select("Which type of project would you like to create?", []string{"Game Client", "Game Server"}, []string{"client", "server"})
		if err != nil {
			return err
		}
	}

	projectName, err := cli.Input("Project name:")
	if err != nil {
		return err
	}

	if _, err := os.Stat(projectName); err == nil {
		return cli.Error("Project '%s' already exists.", projectName)
	}

	err = os.MkdirAll(projectName, 0755)
	if err != nil {
		return err
	}

	switch project {
	case "server":
		err = newServer(projectName)
	case "client":
		err = newClient(projectName)
	default:
		err = cli.Error("Unknown project type: %s", project)
	}

	if err != nil {
		os.RemoveAll(projectName)
		return err
	}

	cli.Success("Successfully created project in '%s/'.", projectName)
	return nil
}

func newServer(projectName string) error {
	return cli.Error("Not implemented.")
}

func newClient(projectName string) error {
	url, err := cli.Input("Enter the URL of the game server:")
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
		language, err = cli.Select("In which language do you want to write your project?", []string{"Go"}, []string{"go"})
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
		cli.Error("Failed to generate event definitions: %s", err)
	}

	switch language {
	case "go":
		err = newClientGo(projectName, url, cgVersion)
	default:
		return cli.Error("Unsupported language: %s", language)
	}
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
	if err != nil || res.StatusCode != http.StatusOK || !external.HasContentType(res.Header, "application/json") {
		return "", "", cli.Error("Couldn't access /info endpoint.")
	}
	defer res.Body.Close()

	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", "", cli.Error("Couldn't decode /info data.")
	}

	return data.Name, data.CGVersion, nil
}

func getCGEVersion(baseURL string) (string, error) {
	res, err := http.Get(baseURL + "/events")
	if err != nil || res.StatusCode != http.StatusOK || !external.HasContentType(res.Header, "text/plain") {
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
