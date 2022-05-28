package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "embed"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/external"
)

//go:embed templates/main.go.tmpl
var goMainTemplate string

func newClientGo(projectName, serverURL, cgVersion string) error {
	err := createGoTemplate(projectName, serverURL)
	if err != nil {
		return err
	}

	err = installGoLibrary(projectName, cgVersion)
	if err != nil {
		return err
	}

	cli.Begin("Cleaning up...")
	err = external.ExecuteInDir(projectName, "go", "mod", "tidy")
	if err != nil {
		return err
	}

	if !external.IsInstalled("goimports") {
		cli.Warn("Failed to add import statements: 'goimports' is not installed!")
		return nil
	}

	external.ExecuteInDir(projectName, "goimports", "-w", "main.go")
	cli.Finish()

	return nil
}

func createGoTemplate(projectName, serverURL string) error {
	module, err := cli.Input("Project module path:")
	if err != nil {
		return err
	}

	cli.Begin("Creating project template...")
	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "init", module)
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	tmpl, err := template.New("main.go").Parse(goMainTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(projectName, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	type data struct {
		URL string
	}

	err = tmpl.Execute(file, data{
		URL: serverURL,
	})
	cli.Finish()
	return err
}

func installGoLibrary(projectName, cgVersion string) error {
	cli.Begin("Fetching correct client library version...")

	clientVersion := external.ClientVersionFromCGVersion("code-game-project", "go-client", cgVersion)

	if clientVersion == "latest" {
		out, err := external.ExecuteInDirHidden(projectName, "go", "get")
		if err != nil && out != "" {
			cli.Error(out)
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
	cli.Finish()

	cli.Begin("Installing dependencies...")
	out, err := external.ExecuteInDirHidden(projectName, "go", "get", path)
	if err != nil {
		cli.Error(out)
		return err
	}
	cli.Finish()
	return nil
}
