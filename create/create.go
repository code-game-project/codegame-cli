package create

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	_ "embed"

	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/exec"
	"github.com/code-game-project/cli-utils/feedback"
	"github.com/code-game-project/cli-utils/modules"
	"github.com/code-game-project/cli-utils/request"
	"github.com/code-game-project/cli-utils/server"
	"github.com/code-game-project/cli-utils/templates"
)

var projectNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-]*$`)

func Create() error {
	projectType := cli.Select("Project type:", []string{"Game client", "Game server"})
	projectName := cli.Input("Project name:", true, "", cli.Regexp(projectNameRegexp, "Project name must only contain 'a'-'z','A'-'Z','0'-'9','-','_'."))

	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("project '%s' already exists", projectName)
	}

	err := os.MkdirAll(projectName, 0o755)
	if err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}
	err = os.Chdir(projectName)
	if err != nil {
		deleteCurrentDir()
		return fmt.Errorf("chdir to project directory: %w", err)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)
	go func() {
		_, ok := <-c
		if ok {
			deleteCurrentDir()
			os.Exit(0)
		}
	}()

	switch projectType {
	case 0:
		err = createClient()
	case 1:
		err = createServer(projectName)
	}
	if err != nil {
		deleteCurrentDir()
		return err
	}

	git()
	readme(projectName)
	license()

	cli.PrintColor(cli.GreenBold, "Successfully created project in '%s/'.", projectName)
	return nil
}

func createClient() error {
	availableLanguages := modules.AvailableLanguages()
	languages := make(map[string]string)
	for name, info := range availableLanguages {
		if info.SupportsClient {
			languages[info.DisplayName] = name
		}
	}
	if len(languages) == 0 {
		return fmt.Errorf("no available language modules")
	}
	lang := cli.SelectString("Language:", languages)
	mod, err := modules.LoadModule(lang)
	if err != nil {
		return fmt.Errorf("load %s module: %w", lang, err)
	}

	var info server.GameInfo
	url := request.TrimURL(cli.Input("Game URL:", true, "", func(input interface{}) error {
		info, err = server.FetchGameInfo(request.TrimURL(input.(string)))
		if err != nil {
			return errors.New("The URL does not point to a valid CodeGame game server or is not reachable")
		}
		return nil
	}))

	return mod.ExecCreateClient(info.Name, url, lang, info.CGVersion)
}

func createServer(projectName string) error {
	return nil
}

func git() {
	if !exec.IsInstalled("git") || !cli.YesNo("Initialize git?", true) {
		os.Remove(".gitignore")
		return
	}

	err := exec.ExecuteDimmed("git", "init")
	if err != nil {
		os.Remove(".gitignore")
		feedback.Error("codegame-cli", "Failed to initialize Git: %s", err)
		return
	}
}

func readme(projectName string) {
	yes := cli.YesNo("Create README?", true)
	if !yes {
		return
	}

	fileContent := fmt.Sprintf("# %s", projectName)

	err := os.WriteFile("README.md", []byte(fileContent), 0o644)
	if err != nil {
		feedback.Error("codegame-cli", "Failed to create README.md: %s", err)
		return
	}
}

//go:embed licenses/MIT.tmpl
var licenseMIT string

//go:embed licenses/MIT_README.tmpl
var licenseReadmeMIT string

//go:embed licenses/GPL.tmpl
var licenseGPL string

//go:embed licenses/GPL_README.tmpl
var licenseReadmeGPL string

//go:embed licenses/AGPL.tmpl
var licenseAGPL string

//go:embed licenses/AGPL_README.tmpl
var licenseReadmeAGPL string

//go:embed licenses/Apache.tmpl
var licenseApache string

//go:embed licenses/Apache_README.tmpl
var licenseReadmeApache string

func license() {
	index := cli.Select("License:", []string{"None", "MIT", "GPLv3", "AGPL", "Apache 2.0"})

	var licenseTemplate string
	var licenseReadmeTemplate string
	switch index {
	case 0:
		return
	case 1:
		licenseTemplate = licenseMIT
		licenseReadmeTemplate = licenseReadmeMIT
	case 2:
		licenseTemplate = licenseGPL
		licenseReadmeTemplate = licenseReadmeGPL
	case 3:
		licenseTemplate = licenseAGPL
		licenseReadmeTemplate = licenseReadmeAGPL
	case 4:
		licenseTemplate = licenseApache
		licenseReadmeTemplate = licenseReadmeApache
	default:
		panic("unknown license")
	}

	username, err := templates.GetUsername()
	if err != nil {
		username = "<your-name>"
		feedback.Warn("codegame-cli", "Make sure to replace <your-name> with your name in LICENSE and README.md")
	}

	err = writeLicense(licenseTemplate, username, time.Now().Year())
	if err != nil {
		os.Remove("LICENSE")
		feedback.Error("codegame-cli", "Failed to create LICENSE: %s", err)
		return
	}

	if _, err := os.Stat("README.md"); err == nil {
		err = writeReadmeLicense(licenseReadmeTemplate, username, time.Now().Year())
		if err != nil {
			feedback.Error("codegame-cli", "Failed to append license information to README.md: %s", err)
			return
		}
	}
}

func writeLicense(templateText, username string, year int) error {
	type data struct {
		Year     int
		Username string
	}
	tmpl, err := template.New("LICENSE").Parse(templateText)
	if err != nil {
		return err
	}

	file, err := os.Create("LICENSE")
	if err != nil {
		return fmt.Errorf("create LICENSE: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data{
		Year:     year,
		Username: username,
	})
}

func writeReadmeLicense(templateText, username string, year int) error {
	readme, err := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0o755)
	if err != nil {
		return fmt.Errorf("open README.md: %w", err)
	}
	defer readme.Close()

	text := "\n\n## License\n\n"
	readme.WriteString(text)

	tmpl, err := template.New("README_License").Parse(templateText)
	if err != nil {
		return err
	}

	type data struct {
		Year     int
		Username string
	}
	return tmpl.Execute(readme, data{
		Year:     year,
		Username: username,
	})
}

func deleteCurrentDir() {
	workingDir, err := os.Getwd()
	if err != nil {
		return
	}

	name := filepath.Base(workingDir)

	os.Chdir("..")

	os.RemoveAll(name)
}
