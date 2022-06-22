package commands

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/util/cgfile"
	"github.com/code-game-project/codegame-cli/util/cggenevents"
	"github.com/code-game-project/codegame-cli/util/exec"
	"github.com/code-game-project/codegame-cli/util/external"
	"github.com/code-game-project/codegame-cli/util/modules"
	"github.com/code-game-project/codegame-cli/util/semver"
)

//go:embed templates/events.cge.tmpl
var eventsCGETemplate string

var projectNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-]*$`)

func New() error {
	var project string
	if len(os.Args) >= 3 {
		project = strings.ToLower(os.Args[2])
	} else {
		var err error
		project, err = cli.SelectString("Project type:", []string{"Game Client", "Game Server"}, []string{"client", "server"})
		if err != nil {
			return err
		}
	}

	projectName, err := cli.Input("Project name:", cli.Regexp(projectNameRegexp, "Project name must only contain 'a'-'z','A'-'Z','0'-'9','-','_'."))
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

	cli.PrintColor(cli.GreenBold, "Successfully created project in '%s/'.", projectName)
	return nil
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
		return cli.Error("Failed to create .codegame.json: %s", err)
	}

	switch language {
	case "go":
		err = modules.Execute("go", "latest", "server", "new", "server")
	case "js":
		err = modules.Execute("js", "latest", "server", "new", "server")
	case "ts":
		err = modules.Execute("js", "latest", "server", "new", "server", "--typescript")
	default:
		return cli.Error("'new server' is not supported for '%s'", language)
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
	url = baseURL(url)
	name, cgVersion, err := getCodeGameInfo(url)
	if err != nil {
		return err
	}
	cgeVersion, err := cggenevents.GetCGEVersion(url)
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
		return cli.Error(err.Error())
	}

	file := cgfile.CodeGameFileData{
		Game: name,
		Type: "client",
		Lang: language,
		URL:  trimURL(url),
	}
	err = file.Write("")
	if err != nil {
		return cli.Error("Failed to create .codegame.json: %s", err)
	}

	switch language {
	case "go":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "go-client", cgVersion)
		err = modules.Execute("go", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url), fmt.Sprintf("--generate-wrappers=%t", cgeMajor > 0 || cgeMinor >= 3))
	case "js":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", cgVersion)
		err = modules.Execute("js", libraryVersion, "client", "new", "client", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url))
	case "ts":
		libraryVersion := external.LibraryVersionFromCGVersion("code-game-project", "javascript-client", cgVersion)
		err = modules.Execute("js", libraryVersion, "client", "new", "client", "--typescript", "--library-version="+libraryVersion, "--game-name="+name, "--url="+trimURL(url))
	default:
		return cli.Error("'new client' is not supported for '%s'", language)
	}
	if err != nil {
		return err
	}

	eventsOutput := "."
	if language == "go" {
		eventsOutput = strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", "")
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

func getCodeGameInfo(baseURL string) (name string, cgVersion string, err error) {
	type response struct {
		Name      string `json:"name"`
		CGVersion string `json:"cg_version"`
		Version   string `json:"version"`
	}
	url := baseURL + "/info"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return "", "", cli.Error("Couldn't access %s.", url)
	}
	if !external.HasContentType(res.Header, "application/json") {
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
