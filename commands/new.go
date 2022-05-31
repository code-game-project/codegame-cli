package commands

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/external"
	"github.com/ogier/pflag"
)

func New() error {
	var project string
	if pflag.NArg() >= 2 {
		project = strings.ToLower(pflag.Arg(1))
	} else {
		var err error
		project, err = cli.Select("Which type of project would you like to create?", []string{"Game Client", "Game Server"}, []string{"client", "server"})
		if err != nil {
			return err
		}
	}

	projectName, err := cli.Input("Project name:")
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

	switch project {
	case "server":
		err = newServer(projectName)
	case "client":
		err = newClient(projectName)
	default:
		err = cli.Error("Unknown project type: %s", project)
	}

	if err != nil {
		os.RemoveAll(projectName)
		return err
	}

	err = git(projectName)
	if err != nil {
		return err
	}
	err = readme(projectName)
	if err != nil {
		return err
	}
	err = license(projectName)
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
		language, err = cli.Select("Language:", []string{"Go"}, []string{"go"})
		if err != nil {
			return err
		}
	}

	var err error
	switch language {
	case "go":
		err = newGoServer(projectName)
	default:
		return cli.Error("Unsupported language: %s", language)
	}
	return err
}

func newClient(projectName string) error {
	url, err := cli.Input("Game server URL:")
	if err != nil {
		return err
	}
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "ws://") {
		url = strings.TrimPrefix(url, "ws://")
	} else if strings.HasPrefix(url, "wss://") {
		url = strings.TrimPrefix(url, "wss://")
	}
	url = strings.TrimSuffix(url, "/")
	ssl := isSSL(url)
	name, cgVersion, err := getCodeGameInfo(baseURL(url, ssl))
	if err != nil {
		return err
	}
	cgeVersion, err := getCGEVersion(baseURL(url, ssl))
	if err != nil {
		return err
	}

	var language string
	if pflag.NArg() >= 3 {
		language = strings.ToLower(pflag.Arg(2))
	} else {
		var err error
		language, err = cli.Select("Language", []string{"Go"}, []string{"go"})
		if err != nil {
			return err
		}
	}

	eventsOutput := projectName
	if language == "go" {
		eventsOutput = filepath.Join(projectName, strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", ""))
	}

	err = external.CGGenEvents(eventsOutput, baseURL(url, ssl), cgeVersion, language)
	if err != nil {
		cli.Error("Failed to generate event definitions: %s", err)
	}

	switch language {
	case "go":
		err = newGoClient(projectName, url, cgVersion)
	default:
		return cli.Error("Unsupported language: %s", language)
	}
	return err
}

func git(projectName string) error {
	if !external.IsInstalled("git") {
		return nil
	}

	yes, err := cli.YesNo("Initialize git?", true)
	if err != nil {
		os.RemoveAll(projectName)
		return err
	}
	if !yes {
		return nil
	}
	out, err := external.ExecuteInDirHidden(projectName, "git", "init")
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
		os.RemoveAll(projectName)
		return err
	}
	if !yes {
		return nil
	}

	fileContent := fmt.Sprintf("# %s", projectName)

	err = os.WriteFile(filepath.Join(projectName, "README.md"), []byte(fileContent), 0644)
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

func license(projectName string) error {
	license, err := cli.Select("License", []string{"None", "MIT", "GPLv3", "AGPL", "Apache 2.0"}, []string{"none", "MIT", "GPL", "AGPL", "Apache"})
	if err != nil {
		os.RemoveAll(projectName)
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

	err = writeLicense(licenseTemplate, projectName, external.GetUsername(), time.Now().Year())
	if err != nil {
		os.Remove(filepath.Join(projectName, "LICENSE"))
		return err
	}

	if _, err := os.Stat("README.md"); err == nil {
		err = writeReadmeLicense(licenseReadmeTemplate, projectName, external.GetUsername(), time.Now().Year())
		if err != nil {
			return err
		}
	}

	return nil
}

func writeLicense(templateText, projectName, username string, year int) error {
	type data struct {
		Year     int
		Username string
	}
	tmpl, err := template.New("LICENSE").Parse(templateText)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(projectName, "LICENSE"))
	if err != nil {
		return cli.Error("Failed to create LICENSE file!")
	}
	defer file.Close()

	return tmpl.Execute(file, data{
		Year:     year,
		Username: username,
	})
}

func writeReadmeLicense(templateText, projectName, username string, year int) error {
	readme, err := os.OpenFile(filepath.Join(projectName, "README.md"), os.O_APPEND|os.O_WRONLY, 0755)
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

func getCGEVersion(baseURL string) (string, error) {
	res, err := http.Get(baseURL + "/events")
	if err != nil || res.StatusCode != http.StatusOK || !external.HasContentType(res.Header, "text/plain") {
		return "", cli.Error("Couldn't access /events endpoint.")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", cli.Error("Couldn't read /events file.")
	}
	return parseCGEVersion([]rune(string(data))), nil
}

func parseCGEVersion(runes []rune) string {
	index := 0
	commentNestingLevel := 0
	for index < len(runes) && (runes[index] == ' ' || runes[index] == '\r' || runes[index] == '\n' || runes[index] == '\t' || (index < len(runes)-1 && runes[index] == '/' && runes[index+1] == '*') || (index < len(runes)-1 && runes[index] == '*' && runes[index+1] == '/') || (index < len(runes)-1 && runes[index] == '/' && runes[index+1] == '/') || commentNestingLevel > 0) {
		if runes[index] == '/' {
			if runes[index+1] == '/' {
				for index < len(runes) && runes[index] != '\n' {
					index++
				}
			} else {
				commentNestingLevel++
			}
		}
		if runes[index] == '*' {
			commentNestingLevel--
		}
		index++
	}

	words := strings.Fields(string(runes[index:]))
	for i, w := range words {
		if w == "version" && i < len(words)-1 {
			return words[i+1]
		}
	}

	return ""
}

func baseURL(domain string, ssl bool) string {
	if ssl {
		return "https://" + domain
	} else {
		return "http://" + domain
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
