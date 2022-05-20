package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
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

func ExecuteHidden(programName string, args ...string) (string, error) {
	cmd := exec.Command(programName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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

	fmt.Println("Make sure to replace <your-name> with your actual name.")
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
		return fmt.Errorf("unsupported platform")
	}
}

func GithubTagFromMinorVersion(owner, repo, version string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	for _, tag := range data {
		if strings.HasPrefix(tag.Name, "v"+version) {
			return tag.Name, nil
		}
	}
	return "", ErrTagNotFound
}
