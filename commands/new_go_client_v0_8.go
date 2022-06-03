package commands

import (
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/code-game-project/codegame-cli/external"
)

//go:embed templates/go/client/v0.8/main.go.tmpl
var goClientMainTemplatev0_8 string

//go:embed templates/go/client/v0.8/wrappers/main.go.tmpl
var goClientWrapperMainTemplatev0_8 string

//go:embed templates/go/client/v0.8/wrappers/game.go.tmpl
var goClientWrapperGameTemplatev0_8 string

//go:embed templates/go/client/v0.8/wrappers/events.go.tmpl
var goClientWrapperEventsTemplatev0_8 string

func createGoClientTemplatev0_8(projectName, modulePath, gameName, serverURL, libraryURL, cgeVersion string, wrappers bool) error {
	if !wrappers {
		return execGoClientMainTemplatev0_8(projectName, serverURL, libraryURL)
	}

	return execGoClientWrappersv0_8(projectName, modulePath, gameName, serverURL, libraryURL, cgeVersion)
}

func execGoClientMainTemplatev0_8(projectName, serverURL, libraryURL string) error {
	type data struct {
		URL        string
		LibraryURL string
	}

	return execTemplate(goClientMainTemplatev0_8, filepath.Join(projectName, "main.go"), data{
		URL:        serverURL,
		LibraryURL: libraryURL,
	})
}

func execGoClientWrappersv0_8(projectName, modulePath, gameName, serverURL, libraryURL, cgeVersion string) error {
	gamePackageName := strings.ReplaceAll(strings.ReplaceAll(gameName, "-", ""), "_", "")

	gameDir := filepath.Join(projectName, strings.ReplaceAll(strings.ReplaceAll(gameName, "-", ""), "_", ""))

	eventNames, err := external.GetEventNames(baseURL(serverURL, isSSL(serverURL)), cgeVersion)
	if err != nil {
		return err
	}

	type event struct {
		Name       string
		PascalName string
	}

	events := make([]event, len(eventNames))
	for i, e := range eventNames {
		pascal := strings.ReplaceAll(e, "_", " ")
		pascal = strings.Title(pascal)
		pascal = strings.ReplaceAll(pascal, " ", "")
		events[i] = event{
			Name:       e,
			PascalName: pascal,
		}
	}

	data := struct {
		URL         string
		LibraryURL  string
		PackageName string
		ModulePath  string
		Events      []event
	}{
		URL:         serverURL,
		LibraryURL:  libraryURL,
		PackageName: gamePackageName,
		ModulePath:  modulePath,
		Events:      events,
	}

	err = execTemplate(goClientWrapperMainTemplatev0_8, filepath.Join(projectName, "main.go"), data)
	if err != nil {
		return err
	}

	err = execTemplate(goClientWrapperGameTemplatev0_8, filepath.Join(gameDir, "game.go"), data)
	if err != nil {
		return err
	}

	err = execTemplate(goClientWrapperEventsTemplatev0_8, filepath.Join(gameDir, "events.go"), data)
	if err != nil {
		return err
	}

	return nil
}
