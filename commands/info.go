package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/external"
	"github.com/mattn/go-colorable"
	"github.com/ogier/pflag"
)

type gameInfo struct {
	Name          string `json:"name"`
	CGVersion     string `json:"cg_version"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	RepositoryURL string `json:"repository_url"`
}

func Info() error {
	var url string
	if pflag.NArg() >= 2 {
		url = strings.ToLower(pflag.Arg(1))
	} else {
		var err error
		url, err = cli.Input("Game server URL:")
		if err != nil {
			return err
		}
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

	info, err := fetchInfo(url)
	if err != nil {
		return err
	}

	printInfo(info)
	return nil
}

func printInfo(info gameInfo) {
	out := colorable.NewColorableStdout()
	printInfoProperty(out, "Display Name", info.DisplayName, 17)
	printInfoProperty(out, "Name", info.Name, 17)
	printInfoProperty(out, "Description", info.Description, 17)
	printInfoProperty(out, "Version", info.Version, 17)
	printInfoProperty(out, "CodeGame Version", info.CGVersion, 17)
	printInfoProperty(out, "Repository", info.RepositoryURL, 17)
}

func printInfoProperty(out io.Writer, name, value string, labelWidth int) {
	if value == "" {
		return
	}

	label := name + ":"
	if labelWidth-utf8.RuneCountInString(label) > 0 {
		label += strings.Repeat(" ", labelWidth-utf8.RuneCountInString(label))
	}

	fmt.Fprintf(out, "\x1b[36m%s\x1b[0m %s\n", label, value)
}

func fetchInfo(url string) (gameInfo, error) {
	url = baseURL(url, isSSL(url)) + "/info"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return gameInfo{}, cli.Error("Couldn't access %s.", url)
	}
	if !external.HasContentType(res.Header, "application/json") {
		return gameInfo{}, cli.Error("%s doesn't return JSON.", url)
	}
	defer res.Body.Close()

	var data gameInfo
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return gameInfo{}, cli.Error("Couldn't decode /info data.")
	}

	return data, nil
}
