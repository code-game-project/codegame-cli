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
	module, err := cli.Input("Project module path:")
	if err != nil {
		return err
	}

	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "init", module)
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Begin("Installing correct go-client version...")
	libraryURL, libraryTag, err := getGoClientLibraryURL(projectName, cgVersion)
	if err != nil {
		return err
	}

	out, err = external.ExecuteInDirHidden(projectName, "go", "get", fmt.Sprintf("%s@%s", libraryURL, libraryTag))
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}
	cli.Finish()

	cli.Begin("Creating project template...")
	err = createGoClientTemplate(projectName, serverURL, libraryURL)
	if err != nil {
		return err
	}
	cli.Finish()

	cli.Begin("Installing dependencies...")

	out, err = external.ExecuteInDirHidden(projectName, "go", "mod", "tidy")
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Finish()

	cli.Begin("Organizing imports...")

	if !external.IsInstalled("goimports") {
		cli.Warn("Failed to organize import statements: 'goimports' is not installed!")
		return nil
	}
	external.ExecuteInDir(projectName, "goimports", "-w", "main.go")

	cli.Finish()

	return nil
}

func createGoClientTemplate(projectName, serverURL, libraryURL string) error {
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
		URL        string
		LibraryURL string
	}

	return tmpl.Execute(file, data{
		URL:        serverURL,
		LibraryURL: libraryURL,
	})
}

func getGoClientLibraryURL(projectName, cgVersion string) (url string, tag string, err error) {
	clientVersion := external.ClientVersionFromCGVersion("code-game-project", "go-client", cgVersion)

	if clientVersion == "latest" {
		var err error
		clientVersion, err = external.LatestGithubTag("code-game-project", "go-client")
		if err != nil {
			return "", "", err
		}
		clientVersion = strings.TrimPrefix(strings.Join(strings.Split(clientVersion, ".")[:2], "."), "v")
	}

	majorVersion := strings.Split(clientVersion, ".")[0]
	tag, err = external.GithubTagFromVersion("code-game-project", "go-client", clientVersion)
	if err != nil {
		return "", "", err
	}
	path := "github.com/code-game-project/go-client/cg"
	if majorVersion != "0" && majorVersion != "1" {
		path = fmt.Sprintf("github.com/code-game-project/go-client/v%s/cg", majorVersion)
	}

	return path, tag, nil
}

func newGoServer(projectName string) error {
	module, err := cli.Input("Project module path:")
	if err != nil {
		return err
	}

	out, err := external.ExecuteInDirHidden(projectName, "go", "mod", "init", module)
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Begin("Fetching latest version numbers...")
	cgeVersion, err := external.LatestCGEVersion()
	if err != nil {
		return err
	}

	libraryURL, err := getServerLibraryURL()
	if err != nil {
		return err
	}
	cli.Finish()

	cli.Begin("Creating project template...")
	err = createGoServerTemplate(projectName, module, cgeVersion, libraryURL)
	if err != nil {
		return err
	}
	cli.Finish()

	cli.Begin("Installing dependencies...")

	out, err = external.ExecuteInDirHidden(projectName, "go", "mod", "tidy")
	if err != nil {
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	cli.Finish()

	cli.Begin("Organizing imports...")

	if !external.IsInstalled("goimports") {
		cli.Warn("Failed to organize import statements: 'goimports' is not installed!")
		return nil
	}
	external.ExecuteInDirHidden(projectName, "goimports", "-w", "main.go")
	packageDir := strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", "")
	external.ExecuteInDirHidden(projectName, "goimports", "-w", filepath.Join(packageDir, "game.go"))

	cli.Finish()

	return nil
}

func createGoServerTemplate(projectName, module, cgeVersion, libraryURL string) error {
	err := executeGoServerTemplate(goServerMainTemplate, "main.go", projectName, cgeVersion, libraryURL, module)
	if err != nil {
		return err
	}

	err = executeGoServerTemplate(goServerCGETemplate, "events.cge", projectName, cgeVersion, libraryURL, module)
	if err != nil {
		return err
	}

	packageName := strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", "")

	err = executeGoServerTemplate(goServerGameTemplate, filepath.Join(packageName, "game.go"), projectName, cgeVersion, libraryURL, module)
	if err != nil {
		return err
	}

	return executeGoServerTemplate(goServerEventsTemplate, filepath.Join(packageName, "events.go"), projectName, cgeVersion, libraryURL, module)
}

func executeGoServerTemplate(templateText, fileName, projectName, cgeVersion, libraryURL, modulePath string) error {
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
		LibraryURL    string
		ModulePath    string
	}

	return tmpl.Execute(file, data{
		Name:          projectName,
		PackageName:   strings.ReplaceAll(strings.ReplaceAll(projectName, "_", ""), "-", ""),
		SnakeCaseName: strings.ReplaceAll(projectName, "-", "_"),
		CGEVersion:    cgeVersion,
		LibraryURL:    libraryURL,
		ModulePath:    modulePath,
	})
}

func getServerLibraryURL() (string, error) {
	tag, err := external.LatestGithubTag("code-game-project", "go-server")
	if err != nil {
		return "", err
	}
	majorVersion := strings.TrimPrefix(strings.Split(tag, ".")[0], "v")

	path := "github.com/code-game-project/go-server/cg"
	if majorVersion != "0" && majorVersion != "1" {
		path = fmt.Sprintf("github.com/code-game-project/go-server/v%s/cg", majorVersion)
	}
	return path, nil
}
