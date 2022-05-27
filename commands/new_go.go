package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-game-project/codegame-cli/external"
	"github.com/code-game-project/codegame-cli/input"
	"github.com/code-game-project/codegame-cli/templates"
)

func newClientGo(projectName, serverURL, cgVersion string) error {
	err := createGoTemplate(projectName, serverURL)
	if err != nil {
		return err
	}

	err = installGoLibrary(projectName, cgVersion)
	if err != nil {
		return err
	}

	err = external.ExecuteInDir(projectName, "go", "mod", "tidy")
	if err != nil {
		return err
	}

	return external.ExecuteInDir(projectName, "goimports", "-w", "main.go")
}

func createGoTemplate(projectName, serverURL string) error {
	module, err := input.Input("Project module path:")
	if err != nil {
		return err
	}

	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "init", module)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return os.WriteFile(filepath.Join(projectName, "main.go"), []byte(strings.ReplaceAll(fmt.Sprintf(templates.Go, serverURL), "$", "%")), 0644)
}

func installGoLibrary(projectName, cgVersion string) error {
	clientVersion := external.ClientVersionFromCGVersion("code-game-project", "go-client", cgVersion)

	if clientVersion == "latest" {
		out, err := external.ExecuteInDirHidden(projectName, "go", "get")
		if err != nil {
			fmt.Println(out)
		}
		return err
	}

	majorVersion := strings.Split(clientVersion, ".")[0]
	tag, err := external.GithubTagFromVersion("code-game-project", "go-client", clientVersion)
	if err != nil {
		return err
	}
	path := "github.com/code-game-project/go-client/cg"
	if majorVersion != "0" && majorVersion != "1" {
		path = fmt.Sprintf("github.com/code-game-project/go-client/v%s/cg", majorVersion)
	}
	path += "@" + tag

	out, err := external.ExecuteInDirHidden(projectName, "go", "get", path)
	if err != nil {
		fmt.Println(out)
	}
	return err
}
