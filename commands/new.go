package commands

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/util"
	"github.com/ogier/pflag"
)

//go:embed templates/events.cge.tmpl
var eventsCGETemplate string

func New() error {
	var project string
	if pflag.NArg() >= 2 {
		project = strings.ToLower(pflag.Arg(1))
	} else {
		var err error
		project, err = cli.Select("Project type:", []string{"Game Client", "Game Server"}, []string{"client", "server"})
		if err != nil {
			return err
		}
	}

	projectName, err := cli.InputAlphanum("Project name:")
	if err != nil {
		return err
	}

	if _, err := os.Stat(projectName); err == nil {
		return cli.Error("Project '%s' already exists.", projectName)
	}

	err = os.MkdirAll(projectName, 0755)
	if err != nil {
		return err
	}
	err = os.Chdir(projectName)
	if err != nil {
		return err
	}

	switch project {
	case "server":
		err = newServer(projectName)
	case "client":
		err = newClient()
	default:
		err = cli.Error("Unknown project type: %s", project)
	}

	if err != nil {
		deleteCurrentDir()
		return err
	}

	err = git()
	if err != nil {
		return err
	}
	err = readme(projectName)
	if err != nil {
		return err
	}
	err = license()
	if err != nil {
		return err
	}

	cli.Success("Successfully created project in '%s/'.", projectName)
	return nil
}

func newServer(projectName string) error {
	var language string
	if pflag.NArg() >= 3 {
		language = strings.ToLower(pflag.Arg(2))
	} else {
		var err error
		language, err = cli.Select("Language:", []string{"Go", "JavaScript", "TypeScript"}, []string{"go", "js", "ts"})
		if err != nil {
			return err
		}
	}

	var err error
	switch language {
	case "go":
		err = util.ExecuteModule("go", "latest", "server", "new", "server")
	case "js":
		err = util.ExecuteModule("js", "latest", "server", "new", "server")
	case "ts":
		err = util.ExecuteModule("js", "latest", "server", "new", "server", "--typescript")
	default:
		return cli.Error("Unsupported language: %s", language)
	}
	if err != nil {
		return err
	}

	cgeVersion, err := util.LatestCGEVersion()
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

	file, err := os.Create("events.cge")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data{
		SnakeCaseName: strings.ReplaceAll(projectName, "-", "_"),
		CGEVersion:    cgeVersion,
	})
}

func newClient() error {
	url, err := cli.Input("Game server URL:")
	if err != nil {
		return err
	}
	url = baseURL(url)
	name, cgVersion, err := getCodeGameInfo(url)
	if err != nil {
		return err
	}
	cgeVersion, err := util.GetCGEVersion(url)
	if err != nil {
		return err
	}

	var language string
	if pflag.NArg() >= 3 {
		language = strings.ToLower(pflag.Arg(2))
	} else {
		var err error
		language, err = cli.Select("Language:", []string{"Go", "JavaScript", "TypeScript"}, []string{"go", "js", "ts"})
		if err != nil {
			return err
		}
	}

	cgeMajor, cgeMinor, _, err := util.ParseVersion(cgeVersion)
	if err != nil {
		return cli.Error(err.Error())
	}

	switch language {
	case "go":
		libraryVersion := util.LibraryVersionFromCGVersion("code-game-project", "go-client", cgVersion)
		err = util.ExecuteModule("go", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url), fmt.Sprintf("--generate-wrappers=%t", cgeMajor > 0 || cgeMinor >= 3))
	case "js":
		libraryVersion := util.LibraryVersionFromCGVersion("code-game-project", "javascript-client", cgVersion)
		err = util.ExecuteModule("js", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url))
	case "ts":
		libraryVersion := util.LibraryVersionFromCGVersion("code-game-project", "javascript-client", cgVersion)
		err = util.ExecuteModule("js", libraryVersion, "client", "new", "client", "--typescript", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url))
	default:
		return cli.Error("Unsupported language: %s", language)
	}
	if err != nil {
		return err
	}

	eventsOutput := "."
	if language == "go" {
		eventsOutput = strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", "")
	}

	if language == "go" || language == "ts" {
		err = util.CGGenEvents(eventsOutput, url, cgeVersion, language)
		if err != nil {
			return err
		}
	}

	return nil
}

func git() error {
	if !util.IsInstalled("git") {
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
	out, err := util.Execute(true, "git", "init")
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

	err = os.WriteFile("README.md", []byte(fileContent), 0644)
	if err != nil {
		cli.Error(err.Error())
	}
	return err
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
	license, err := cli.Select("License:", []string{"None", "MIT", "GPLv3", "AGPL", "Apache 2.0"}, []string{"none", "MIT", "GPL", "AGPL", "Apache"})
	if err != nil {
		deleteCurrentDir()
		return err
	}

	var licenseTemplate string
	var licenseReadmeTemplate string
	switch license {
	case "MIT":
		licenseTemplate = licenseMIT
		licenseReadmeTemplate = licenseReadmeMIT
	case "GPL":
		licenseTemplate = licenseGPL
		licenseReadmeTemplate = licenseReadmeGPL
	case "AGPL":
		licenseTemplate = licenseAGPL
		licenseReadmeTemplate = licenseReadmeAGPL
	case "Apache":
		licenseTemplate = licenseApache
		licenseReadmeTemplate = licenseReadmeApache
	case "none":
		return nil
	default:
		return errors.New("Unknown license.")
	}

	err = writeLicense(licenseTemplate, util.GetUsername(), time.Now().Year())
	if err != nil {
		os.Remove("LICENSE")
		return err
	}

	if _, err := os.Stat("README.md"); err == nil {
		err = writeReadmeLicense(licenseReadmeTemplate, util.GetUsername(), time.Now().Year())
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
		return cli.Error("Failed to create LICENSE file!")
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
		return cli.Error("Failed to append license text to README.")
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

func getCodeGameInfo(baseURL string) (string, string, error) {
	type response struct {
		Name      string `json:"name"`
		CGVersion string `json:"cg_version"`
	}
	url := baseURL + "/info"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return "", "", cli.Error("Couldn't access %s.", url)
	}
	if !util.HasContentType(res.Header, "application/json") {
		return "", "", cli.Error("%s doesn't return JSON.", url)
	}
	defer res.Body.Close()

	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", "", cli.Error("Couldn't decode /info data.")
	}

	return data.Name, data.CGVersion, nil
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

func trimURL(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "ws://") {
		url = strings.TrimPrefix(url, "ws://")
	} else if strings.HasPrefix(url, "wss://") {
		url = strings.TrimPrefix(url, "wss://")
	}
	return strings.TrimSuffix(url, "/")
}

// baseURL returns the URL with the correct protocol ('http://' or 'https://')
func baseURL(url string) string {
	url = trimURL(url)
	if isSSL(url) {
		return "https://" + url
	} else {
		return "http://" + url
	}
}

func isSSL(domain string) bool {
	res, err := http.Get("https://" + domain)
	if err == nil {
		res.Body.Close()
		return true
	}
	return false
}
