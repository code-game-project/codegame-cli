package exec

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func IsInstalled(programName string) bool {
	_, err := exec.LookPath(programName)
	return err == nil
}

// Execute the program with args.
// If hidden is set, no output (except errors) will be printed and stdout will not be passed to the program.
// It returns the combined output if hidden is true. Otherwise all output will be printed directly to stdout.
func Execute(hidden bool, programName string, args ...string) (string, error) {
	if _, err := exec.LookPath(programName); err != nil {
		return "", fmt.Errorf("'%s' ist not installed!", programName)
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
			err = fmt.Errorf("'%s' returned with an error:\n%s", programName, outStr)
		} else {
			err = fmt.Errorf("Failed to execute '%s %s'.", programName, strings.Join(args, " "))
		}
	}

	return outStr, err
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
