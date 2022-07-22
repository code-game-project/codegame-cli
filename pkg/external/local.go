package external

import (
	"os/user"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/pkg/exec"
)

// GetUsername tries to determine the name of the current user.
// It looks at the following things in order:
//   1. git config user.name
//   2. currently logged in user of the OS
//   3. returns <your-name> with a note to change it
func GetUsername() string {
	name, err := exec.Execute(true, "git", "config", "user.name")
	if err == nil {
		return strings.TrimSpace(name)
	}

	user, err := user.Current()
	if err == nil {
		return strings.TrimSpace(user.Username)
	}

	cli.PrintColor(cli.Yellow, "Make sure to replace <your-name> with your actual name.")
	return "<your-name>"
}
