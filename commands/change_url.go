package commands

import (
	"os"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/util/cgfile"
)

func ChangeURL() error {
	root, err := cgfile.FindProjectRoot()
	if err != nil {
		return err
	}

	var url string
	if len(os.Args) >= 3 {
		url = strings.ToLower(os.Args[2])
	} else {
		var err error
		url, err = cli.Input("New game URL:")
		if err != nil {
			return err
		}
	}

	config, err := cgfile.LoadCodeGameFile(root)
	if err != nil {
		return err
	}

	if config.Type != "client" {
		return cli.Error("Project is not a client.")
	}

	name, _, err := getCodeGameInfo(baseURL(url))
	if err != nil {
		return err
	}
	if name != config.Game {
		return cli.Error("The URL points to a different game.")
	}

	prevURL := config.URL

	config.URL = url
	err = config.Write(root)
	if err != nil {
		return err
	}

	err = Update()
	if err != nil {
		config.URL = prevURL
		config.Write(root)
		return err
	}
	return nil
}
