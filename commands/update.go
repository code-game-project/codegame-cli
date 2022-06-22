package commands

import (
	"os"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/util/cgfile"
	"github.com/code-game-project/codegame-cli/util/cggenevents"
	"github.com/code-game-project/codegame-cli/util/external"
	"github.com/code-game-project/codegame-cli/util/modules"
)

func Update() error {
	root, err := cgfile.FindProjectRoot()
	if err != nil {
		return err
	}

	err = os.Chdir(root)
	if err != nil {
		return err
	}

	data, err := cgfile.LoadCodeGameFile("")
	if err != nil {
		return cli.Error("Failed to load .codegame.json")
	}

	switch data.Type {
	case "client":
		return updateClient(data)
	case "server":
		return updateServer(data)
	default:
		return cli.Error("Unknown project type: %s", data.Type)
	}

}

func updateClient(config *cgfile.CodeGameFileData) error {
	baseURL := baseURL(config.URL)

	_, cgVersion, err := getCodeGameInfo(baseURL)
	if err != nil {
		return err
	}

	cgeVersion, err := cggenevents.GetCGEVersion(baseURL)
	if err != nil {
		return err
	}

	switch config.Lang {
	case "go":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "go-client", cgVersion)
		err = modules.Execute("go", libraryVersion, "client", "update", "--library-version="+libraryVersion)
	default:
		return cli.Error("'update' is not supported for '%s'", config.Lang)
	}
	if err != nil {
		return err
	}

	eventsOutput := "."
	if config.Lang == "go" {
		eventsOutput = strings.ReplaceAll(strings.ReplaceAll(config.Game, "-", ""), "_", "")
	}

	if config.Lang == "go" || config.Lang == "ts" {
		err = cggenevents.CGGenEvents(cgeVersion, eventsOutput, baseURL, config.Lang)
	}

	return err
}

func updateServer(config *cgfile.CodeGameFileData) error {
	switch config.Lang {
	case "go":
		return modules.Execute("go", "latest", "server", "update")
	default:
		return cli.Error("'update' is not supported for '%s'", config.Lang)
	}
}
