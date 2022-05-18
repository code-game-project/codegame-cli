package external

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

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
