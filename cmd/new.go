package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/cggenevents"
	"github.com/code-game-project/go-utils/exec"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/modules"
	"github.com/code-game-project/go-utils/server"
	"github.com/spf13/cobra"
)

//go:embed templates/events.cge.tmpl
var eventsCGETemplate string

var projectNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-]*$`)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new CodeGame application.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var project string
		if len(args) > 0 {
			project = strings.ToLower(args[0])
		} else {
			var err error
			project, err = cli.SelectString("Project type:", []string{"Game Client", "Game Server"}, []string{"client", "server"})
			abort(err)
		}

		projectName, err := cli.Input("Project name:", cli.Regexp(projectNameRegexp, "Project name must only contain 'a'-'z','A'-'Z','0'-'9','-','_'."))
		abort(err)

		if _, err := os.Stat(projectName); err == nil {
			abort(fmt.Errorf("project '%s' already exists.", projectName))
		}

		err = os.MkdirAll(projectName, 0o755)
		abort(err)
		err = os.Chdir(projectName)
		abort(err)

		switch project {
		case "server":
			err = newServer(projectName)
		case "client":
			err = newClient()
		default:
			err = fmt.Errorf("unknown project type: %s", project)
		}

		if err != nil {
			deleteCurrentDir()
			abort(err)
		}

		err = git()
		abort(err)
		err = readme(projectName)
		abort(err)
		err = license()
		abort(err)

		cli.PrintColor(cli.GreenBold, "Successfully created project in '%s/'.", projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func newServer(projectName string) error {
	language, err := cli.SelectString("Language:", []string{"Go"}, []string{"go"})
	if err != nil {
		return err
	}

	file := cgfile.CodeGameFileData{
		Game: projectName,
		Type: "server",
		Lang: language,
	}
	err = file.Write("")
	if err != nil {
		return fmt.Errorf("Failed to create .codegame.json: %w", err)
	}

	newData := modules.NewServerData{
		Lang:           language,
		LibraryVersion: "latest",
	}

	switch language {
	case "go":
		err = modules.ExecuteNewServer(newData)
	default:
		err = fmt.Errorf("'new server' is not supported for '%s'", language)
	}
	if err != nil {
		return err
	}

	cgeVersion, err := cggenevents.LatestCGEVersion()
	if err != nil {
		return err
	}

	type data struct {
		SnakeCaseName string
		CGEVersion    string
	}

	tmpl, err := template.New("events.cge").Parse(eventsCGETemplate)
	if err != nil {
		return err
	}

	eventsFile, err := os.Create("events.cge")
	if err != nil {
		return err
	}
	defer eventsFile.Close()

	return tmpl.Execute(eventsFile, data{
		SnakeCaseName: strings.ReplaceAll(projectName, "-", "_"),
		CGEVersion:    cgeVersion,
	})
}

func newClient() error {
	url, err := cli.Input("Game server URL:")
	if err != nil {
		return err
	}
	url = external.TrimURL(url)
	api, err := server.NewAPI(url)
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

	language, err := cli.SelectString("Language:", []string{"C#", "Go", "Java", "JavaScript", "TypeScript"}, []string{"cs", "go", "java", "js", "ts"})
	if err != nil {
		return err
	}

	file := &cgfile.CodeGameFileData{
		Game:        info.Name,
		GameVersion: info.Version,
		Type:        "client",
		Lang:        language,
		URL:         url,
	}
	err = file.Write("")
	if err != nil {
		return fmt.Errorf("Failed to create .codegame.json: %s", err)
	}

	newData := modules.NewClientData{
		Lang: language,
		Name: info.Name,
		URL:  url,
	}

	switch language {
	case "cs":
		newData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "csharp-client", info.CGVersion)
		err = modules.ExecuteNewClient(newData)
	case "go":
		newData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "go-client", info.CGVersion)
		err = modules.ExecuteNewClient(newData)
	case "java":
		newData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "java-client", info.CGVersion)
		err = modules.ExecuteNewClient(newData)
	case "js", "ts":
		newData.LibraryVersion = external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", info.CGVersion)
		err = modules.ExecuteNewClient(newData)
	default:
		err = fmt.Errorf("'new client' is not supported for '%s'", language)
	}
	if err != nil {
		return err
	}

	file, err = cgfile.LoadCodeGameFile("")
	if err != nil {
		return fmt.Errorf("Failed to open .codegame.json: %w", err)
	}

	if language == "cs" || language == "go" || language == "java" || language == "ts" {
		eventsOutput := info.Name
		switch language {
		case "cs":
			eventsOutput = strings.ReplaceAll(strings.Title(strings.ReplaceAll(strings.ReplaceAll(eventsOutput, "_", " "), "-", " ")), " ", "")
		case "go":
			eventsOutput = strings.ReplaceAll(strings.ReplaceAll(eventsOutput, "-", ""), "_", "")
		case "java":
			packageConf, ok := file.LangConfig["package"]
			if !ok {
				return errors.New("Missing language config field `package` in .codegame.json!")
			}
			packageName := packageConf.(string)
			if packageConf == "" {
				return errors.New("Empty language config field `package` in .codegame.json!")
			}
			gameDir := filepath.Join("src", "main", "java")
			pkgDir := filepath.Join(strings.Split(packageName, ".")...)
			eventsOutput = filepath.Join(gameDir, pkgDir, strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(eventsOutput), "_", ""), "-", ""))
		case "ts":
			eventsOutput = filepath.Join("src", eventsOutput)
		}
		err = cggenevents.CGGenEvents(cgeVersion, eventsOutput, external.BaseURL("http", external.IsTLS(url), url), language)
		if err != nil {
			return err
		}
	}

	return nil
}

func git() error {
	if !exec.IsInstalled("git") {
		os.Remove(".gitignore")
		return nil
	}

	yes, err := cli.YesNo("Initialize git?", true)
	if err != nil {
		deleteCurrentDir()
		return err
	}
	if !yes {
		os.Remove(".gitignore")
		return nil
	}
	out, err := exec.Execute(true, "git", "init")
	if err != nil {
		os.Remove(".gitignore")
		if out != "" {
			cli.Error(out)
		}
		return err
	}

	return nil
}

func readme(projectName string) error {
	yes, err := cli.YesNo("Create README?", true)
	if err != nil {
		deleteCurrentDir()
		return err
	}
	if !yes {
		return nil
	}

	fileContent := fmt.Sprintf("# %s", projectName)

	return os.WriteFile("README.md", []byte(fileContent), 0o644)
}

//go:embed templates/licenses/MIT.tmpl
var licenseMIT string

//go:embed templates/licenses/MIT_README.tmpl
var licenseReadmeMIT string

//go:embed templates/licenses/GPL.tmpl
var licenseGPL string

//go:embed templates/licenses/GPL_README.tmpl
var licenseReadmeGPL string

//go:embed templates/licenses/AGPL.tmpl
var licenseAGPL string

//go:embed templates/licenses/AGPL_README.tmpl
var licenseReadmeAGPL string

//go:embed templates/licenses/Apache.tmpl
var licenseApache string

//go:embed templates/licenses/Apache_README.tmpl
var licenseReadmeApache string

func license() error {
	index, err := cli.Select("License:", []string{"None", "MIT", "GPLv3", "AGPL", "Apache 2.0"})
	if err != nil {
		deleteCurrentDir()
		return err
	}

	var licenseTemplate string
	var licenseReadmeTemplate string
	switch index {
	case 0:
		return nil
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
		return errors.New("Unknown license.")
	}

	err = writeLicense(licenseTemplate, external.GetUsername(), time.Now().Year())
	if err != nil {
		os.Remove("LICENSE")
		return err
	}

	if _, err := os.Stat("README.md"); err == nil {
		err = writeReadmeLicense(licenseReadmeTemplate, external.GetUsername(), time.Now().Year())
		if err != nil {
			cli.Error("Failed to write license into README.md")
			return nil
		}
	}

	return nil
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
		return fmt.Errorf("Failed to create LICENSE file!")
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
		return fmt.Errorf("Failed to append license text to README.")
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

func execTemplate(templateText, path string, data any) error {
	err := os.MkdirAll(filepath.Join(filepath.Dir(path)), 0o755)
	if err != nil {
		return err
	}

	tmpl, err := template.New(path).Parse(templateText)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func deleteCurrentDir() {
	workingDir, err := os.Getwd()
	if err != nil {
		cli.Error("Failed to delete created directory: %s", err)
		return
	}

	name := filepath.Base(workingDir)

	os.Chdir("..")

	os.RemoveAll(name)
}
