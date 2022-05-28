package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"github.com/code-game-project/codegame-cli/cli"
)

var ErrTagNotFound = errors.New("tag not found")

func IsInstalled(programName string) bool {
	_, err := exec.LookPath(programName)
	return err == nil
}

func Execute(programName string, args ...string) error {
	cmd := exec.Command(programName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExecuteInDir(workingDir, programName string, args ...string) error {
	cmd := exec.Command(programName, args...)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExecuteHidden(programName string, args ...string) (string, error) {
	cmd := exec.Command(programName, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ExecuteInDirHidden(workingDir, programName string, args ...string) (string, error) {
	cmd := exec.Command(programName, args...)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func GetUsername() string {
	name, err := ExecuteHidden("git", "config", "user.name")
	if err == nil {
		return strings.TrimSpace(name)
	}

	user, err := user.Current()
	if err == nil {
		return strings.TrimSpace(user.Username)
	}

	cli.Info("Make sure to replace <your-name> with your actual name.")
	return "<your-name>"
}

func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("Unsupported platform.")
	}
}

func GithubTagFromVersion(owner, repo, version string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK {
		return "", cli.Error("Couldn't access git tags from 'github.com/%s/%s'.", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", cli.Error("Couldn't decode git tag data.")
	}

	for _, tag := range data {
		if strings.HasPrefix(tag.Name, "v"+version) {
			return tag.Name, nil
		}
	}
	return "", ErrTagNotFound
}

func ClientVersionFromCGVersion(owner, repo, cgVersion string) string {
	res, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/versions.json", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK {
		cli.Warn("Couldn't fetch versions.json. Using latest client library version.")
		return "latest"
	}
	defer res.Body.Close()

	var versions map[string]string
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest client library version.")
		return "latest"
	}

	// check exact match
	if v, ok := versions[cgVersion]; ok {
		return v
	}

	parts := strings.Split(cgVersion, ".")
	if len(parts) < 2 {
		cli.Warn("Invalid versions.json. Using latest client library version.")
		return "latest"
	}
	major := parts[0]

	// get all minor versions of the requested major version
	compatibleMinorVersions := make([]int, 0)
	for v := range versions {
		clientParts := strings.Split(v, ".")
		if len(clientParts) < 2 {
			cli.Warn("Invalid versions.json. Using latest client library version.")
			return "latest"
		}
		clientMajor := clientParts[0]
		if major == clientMajor {
			minor, err := strconv.Atoi(clientParts[1])
			if err != nil {
				cli.Warn("Invalid versions.json. Using latest client library version.")
				return "latest"
			}
			compatibleMinorVersions = append(compatibleMinorVersions, minor)
		}
	}
	if len(compatibleMinorVersions) == 0 {
		cli.Warn("No compatible client library version found. Using latest client library version.")
		return "latest"
	}

	minorStr := parts[1]
	minor, err := strconv.Atoi(minorStr)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest client library version.")
		return "latest"
	}

	// check closest minor version above requested
	closestMinor := -1
	for _, v := range compatibleMinorVersions {
		if v > minor && float64(closestMinor-minor) > float64(v-minor) {
			closestMinor = v
		}
	}
	if closestMinor >= 0 {
		v := fmt.Sprintf("%s.%d", major, minor)
		cli.Warn("No exact version match found. Using client library version %s.", v)
		return v
	}

	// check closest minor version below requested
	closestMinor = math.MaxInt
	for _, v := range compatibleMinorVersions {
		if v < minor && float64(minor-closestMinor) > float64(minor-v) {
			closestMinor = v
		}
	}
	if closestMinor >= 0 {
		v := fmt.Sprintf("%s.%d", major, minor)
		cli.Warn("No exact version match found. Using client library version %s.")
		return v
	}

	cli.Warn("No compatible client library version found. Using latest client library version.")
	return "latest"
}
