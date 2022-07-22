/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	"github.com/code-game-project/codegame-cli/pkg/cgfile"
	"github.com/code-game-project/codegame-cli/pkg/cggenevents"
	"github.com/code-game-project/codegame-cli/pkg/exec"
	"github.com/code-game-project/codegame-cli/pkg/external"
	"github.com/code-game-project/codegame-cli/pkg/modules"
	"github.com/code-game-project/codegame-cli/pkg/semver"
	"github.com/code-game-project/codegame-cli/pkg/server"
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

		err = os.MkdirAll(projectName, 0755)
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
	var language string
	if len(os.Args) >= 4 {
		language = strings.ToLower(os.Args[3])
	} else {
		var err error
		language, err = cli.SelectString("Language:", []string{"Go", "JavaScript", "TypeScript"}, []string{"go", "js", "ts"})
		if err != nil {
			return err
		}
	}

	file := cgfile.CodeGameFileData{
		Game: projectName,
		Type: "server",
		Lang: language,
	}
	err := file.Write("")
	if err != nil {
		return fmt.Errorf("Failed to create .codegame.json: %s", err)
	}

	switch language {
	case "go":
		err = modules.Execute("go", "latest", "server", "new", "server")
	case "js":
		err = modules.Execute("js", "latest", "server", "new", "server")
	case "ts":
		err = modules.Execute("js", "latest", "server", "new", "server", "--typescript")
	default:
		return fmt.Errorf("'new server' is not supported for '%s'", language)
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

	var language string
	if len(os.Args) >= 4 {
		language = strings.ToLower(os.Args[3])
	} else {
		var err error
		language, err = cli.SelectString("Language:", []string{"Go", "JavaScript", "TypeScript"}, []string{"go", "js", "ts"})
		if err != nil {
			return err
		}
	}

	cgeMajor, cgeMinor, _, err := semver.ParseVersion(cgeVersion)
	if err != nil {
		return err
	}

	file := cgfile.CodeGameFileData{
		Game: info.Name,
		Type: "client",
		Lang: language,
		URL:  external.TrimURL(url),
	}
	err = file.Write("")
	if err != nil {
		return fmt.Errorf("Failed to create .codegame.json: %s", err)
	}

	switch language {
	case "go":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "go-client", info.CGVersion)
		err = modules.Execute("go", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+info.Name, "--url="+external.TrimURL(url), fmt.Sprintf("--generate-wrappers=%t", cgeMajor > 0 || cgeMinor >= 3))
	case "js":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", info.CGVersion)
		err = modules.Execute("js", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+info.Name, "--url="+external.TrimURL(url))
	case "ts":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", info.CGVersion)
		err = modules.Execute("js", libraryVersion, "client", "new", "client", "--typescript", "--library-version="+libraryVersion, "--game-name="+info.Name, "--url="+external.TrimURL(url))
	default:
		return fmt.Errorf("'new client' is not supported for '%s'", language)
	}
	if err != nil {
		return err
	}

	eventsOutput := "."
	if language == "go" {
		eventsOutput = strings.ReplaceAll(strings.ReplaceAll(info.Name, "-", ""), "_", "")
	}

	if language == "go" || language == "ts" {
		err = cggenevents.CGGenEvents(cgeVersion, eventsOutput, url, language)
		if err != nil {
			return err
		}
	}

	return nil
}

func git() error {
	if !exec.IsInstalled("git") {
		return nil
	}

	yes, err := cli.YesNo("Initialize git?", true)
	if err != nil {
		deleteCurrentDir()
		return err
	}
	if !yes {
		return nil
	}
	out, err := exec.Execute(true, "git", "init")
	if err != nil {
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

	return os.WriteFile("README.md", []byte(fileContent), 0644)
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
	readme, err := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0755)
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
	err := os.MkdirAll(filepath.Join(filepath.Dir(path)), 0755)
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
