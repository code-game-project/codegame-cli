/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/code-game-project/codegame-cli/pkg/cgfile"
	"github.com/code-game-project/codegame-cli/pkg/cggenevents"
	"github.com/code-game-project/codegame-cli/pkg/external"
	"github.com/code-game-project/codegame-cli/pkg/modules"
	"github.com/code-game-project/codegame-cli/pkg/server"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the current project",
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

	switch config.Lang {
	case "go":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "go-client", info.CGVersion)
		err = modules.Execute("go", libraryVersion, "client", "update", "--library-version="+libraryVersion)
	default:
		return fmt.Errorf("'update' is not supported for '%s'", config.Lang)
	}
	if err != nil {
		return err
	}

	eventsOutput := "."
	if config.Lang == "go" {
		eventsOutput = strings.ReplaceAll(strings.ReplaceAll(config.Game, "-", ""), "_", "")
	}

	if config.Lang == "go" || config.Lang == "ts" {
		err = cggenevents.CGGenEvents(cgeVersion, eventsOutput, api.BaseURL(), config.Lang)
	}

	return err
}

func updateServer(config *cgfile.CodeGameFileData) error {
	switch config.Lang {
	case "go":
		return modules.Execute("go", "latest", "server", "update")
	default:
		return fmt.Errorf("'update' is not supported for '%s'", config.Lang)
	}
}
