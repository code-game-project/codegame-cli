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

//go:embed templates/go/client/main.go.tmpl
var goClientMainTemplate string

//go:embed templates/go/server/main.go.tmpl
var goServerMainTemplate string

//go:embed templates/go/server/game.go.tmpl
var goServerGameTemplate string

//go:embed templates/go/server/events.cge.tmpl
var goServerCGETemplate string

//go:embed templates/go/server/events.go.tmpl
var goServerEventsTemplate string

func newGoClient(projectName, serverURL, cgVersion string) error {
	err := createGoClientTemplate(projectName, serverURL)
	if err != nil {
		return err
	}

	err = installGoClientLibrary(projectName, cgVersion)
	if err != nil {
		return err
	}

	cli.Begin("Cleaning up...")

	if !external.IsInstalled("goimports") {
		cli.Warn("Failed to add import statements: 'goimports' is not installed!")
		return nil
	}

	external.ExecuteInDir(projectName, "goimports", "-w", "main.go")

	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "tidy")
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Finish()

	return nil
}

func createGoClientTemplate(projectName, serverURL string) error {
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

	tmpl, err := template.New("main.go").Parse(goClientMainTemplate)
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

func installGoClientLibrary(projectName, cgVersion string) error {
	cli.Begin("Fetching correct client library version...")

	clientVersion := external.ClientVersionFromCGVersion("code-game-project", "go-client", cgVersion)

	if clientVersion == "latest" {
		var err error
		clientVersion, err = external.LatestGithubTag("code-game-project", "go-client")
		if err != nil {
			return err
		}
		clientVersion = strings.TrimPrefix(strings.Join(strings.Split(clientVersion, ".")[:2], "."), "v")
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

func newGoServer(projectName string) error {
	err := createGoServerTemplate(projectName)
	if err != nil {
		return err
	}

	err = installGoServerLibrary(projectName)
	if err != nil {
		return err
	}

	cli.Begin("Cleaning up...")

	if !external.IsInstalled("goimports") {
		cli.Warn("Failed to add import statements: 'goimports' is not installed!")
		return nil
	}

	external.ExecuteInDir(projectName, "goimports", "-w", "main.go")
	packageDir := strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", "")
	external.ExecuteInDir(projectName, "goimports", "-w", filepath.Join(packageDir, "game.go"))

	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "tidy")
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Finish()

	return nil
}

func createGoServerTemplate(projectName string) error {
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

	cgeVersion, err := external.LatestCGEVersion()
	if err != nil {
		return err
	}

	err = executeGoServerTemplate(goServerMainTemplate, "main.go", projectName, cgeVersion)
	if err != nil {
		return err
	}

	err = executeGoServerTemplate(goServerCGETemplate, "events.cge", projectName, cgeVersion)
	if err != nil {
		return err
	}

	packageName := strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", "")

	err = executeGoServerTemplate(goServerGameTemplate, filepath.Join(packageName, "game.go"), projectName, cgeVersion)
	if err != nil {
		return err
	}

	err = executeGoServerTemplate(goServerEventsTemplate, filepath.Join(packageName, "events.go"), projectName, cgeVersion)
	if err != nil {
		return err
	}

	cli.Finish()
	return nil
}

func executeGoServerTemplate(templateText, fileName, projectName, cgeVersion string) error {
	tmpl, err := template.New(fileName).Parse(templateText)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(projectName, filepath.Dir(fileName)), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(projectName, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	type data struct {
		Name          string
		PackageName   string
		SnakeCaseName string
		CGEVersion    string
	}

	return tmpl.Execute(file, data{
		Name:          projectName,
		PackageName:   strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", ""),
		SnakeCaseName: strings.ReplaceAll(projectName, "-", "_"),
		CGEVersion:    cgeVersion,
	})
}

func installGoServerLibrary(projectName string) error {
	cli.Begin("Fetching latest server library version...")
	tag, err := external.LatestGithubTag("code-game-project", "go-server")
	if err != nil {
		return err
	}
	majorVersion := strings.TrimPrefix(strings.Split(tag, ".")[0], "v")

	path := "github.com/code-game-project/go-server/cg"
	if majorVersion != "0" && majorVersion != "1" {
		path = fmt.Sprintf("github.com/code-game-project/go-server/v%s/cg", majorVersion)
	}
	cli.Finish()

	cli.Begin("Installing dependencies...")

	out, err := external.ExecuteInDirHidden(projectName, "go", "get", path)
	if err != nil {
		cli.Error(out)
		return err
	}

	out, err = external.ExecuteInDirHidden(projectName, "go", "get", "github.com/spf13/pflag")
	if err != nil {
		cli.Error(out)
		return err
	}

	cli.Finish()
	return nil
}
