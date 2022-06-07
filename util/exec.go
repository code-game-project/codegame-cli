package util

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/code-game-project/codegame-cli/cli"
)

var ErrTagNotFound = errors.New("tag not found")

func IsInstalled(programName string) bool {
	_, err := exec.LookPath(programName)
	return err == nil
}

// Execute the program with args.
// If hidden is set, no output (except errors) will be printed and stdout will not be passed to the program.
// It returns the combined output if hidden is true. Otherwise all output will be printed directly to stdout.
func Execute(hidden bool, programName string, args ...string) (string, error) {
	if _, err := exec.LookPath(programName); err != nil {
		cli.Error("'%s' ist not installed!", programName)
		return "", err
	}
	cmd := exec.Command(programName, args...)

	var out []byte
	var err error

	if !hidden {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	} else {
		out, err = cmd.CombinedOutput()
	}

	outStr := string(out)
	if err != nil {
		if outStr != "" {
			err = cli.Error("'%s' returned with an error:\n%s", programName, outStr)
		} else {
			err = cli.Error("Failed to execute '%s'.", programName)
		}
	}

	return outStr, err
}

// GetUsername tries to determine the name of the current user.
// It looks at the following things in order:
//   1. git config user.name
//   2. currently logged in user of the OS
//   3. returns <your-name> with a note to change it
func GetUsername() string {
	name, err := Execute(true, "git", "config", "user.name")
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

// Opens the specified URL in the default browser.
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

// HasContentType returns true if the 'content-type' header includes mimetype.
func HasContentType(h http.Header, mimetype string) bool {
	contentType := h.Get("content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
