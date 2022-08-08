package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/cggenevents"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/modules"
	"github.com/code-game-project/go-utils/server"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the current project.",
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() error {
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
		return fmt.Errorf("Failed to load .codegame.json")
	}
	data.URL = external.TrimURL(data.URL)
	err = data.Write("")
	if err != nil {
		return err
	}

	switch data.Type {
	case "client":
		return updateClient(data)
	case "server":
		return updateServer(data)
	default:
		return fmt.Errorf("Unknown project type: %s", data.Type)
	}
}

func updateClient(config *cgfile.CodeGameFileData) error {
	api, err := server.NewAPI(config.URL)
	if err != nil {
		return err
	}

	info, err := api.FetchGameInfo()
	if err != nil {
		return err
	}

	cge, err := api.GetCGEFile()
	if err != nil {
		return err
	}

	cgeVersion, err := cggenevents.ParseCGEVersion(cge)
	if err != nil {
		return err
	}

	updateData := modules.UpdateData{
		Lang: config.Lang,
	}

	switch config.Lang {
	case "go":
		updateData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "go-client", info.CGVersion)
		err = modules.ExecuteUpdate(updateData, config)
	case "js", "ts":
		updateData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", info.CGVersion)
		err = modules.ExecuteUpdate(updateData, config)
	default:
		err = fmt.Errorf("'update' is not supported for '%s'", config.Lang)
	}
	if err != nil {
		return err
	}

	if config.Lang == "go" || config.Lang == "ts" {
		eventsOutput := config.Game
		if config.Lang == "go" {
			eventsOutput = strings.ReplaceAll(strings.ReplaceAll(eventsOutput, "-", ""), "_", "")
		} else if config.Lang == "ts" {
			eventsOutput = filepath.Join("src", eventsOutput)
		}
		err = cggenevents.CGGenEvents(cgeVersion, eventsOutput, api.BaseURL(), config.Lang)
	}

	return err
}

func updateServer(config *cgfile.CodeGameFileData) error {
	updateData := modules.UpdateData{
		Lang:           config.Lang,
		LibraryVersion: "latest",
	}

	switch config.Lang {
	case "go":
		return modules.ExecuteUpdate(updateData, config)
	default:
		return fmt.Errorf("'update' is not supported for '%s'", config.Lang)
	}
}
