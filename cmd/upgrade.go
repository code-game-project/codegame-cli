package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Bananenpro/cli"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update codegame-cli to the latest version.",
	Run: func(cmd *cobra.Command, args []string) {
		if rootCmd.Version == "dev" {
			cli.Error("Cannot update dev version.")
			os.Exit(1)
		}

		if !versionCheck(false, true) {
			cli.Success("codegame-cli is already up-to-date.")
			return
		}

		switch runtime.GOOS {
		case "windows":
			upgradeWindows()
		case "darwin", "linux":
			upgradeUnix()
		default:
			cli.Error("Automatic updates are not supported for your operating system.")
			os.Exit(1)
		}
	},
}

func upgradeWindows() {
	fmt.Println("Downloading latest installer...")
	cmd := exec.Command("Powershell.exe", "-Command", "iwr -useb https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.ps1 | iex")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		cli.Error("Failed to update codegame-cli.")
		os.Exit(1)
	}
}

func upgradeUnix() {
	fmt.Println("Downloading latest installer...")

	var installCmd string
	if _, err := exec.LookPath("wget"); err == nil {
		installCmd = "wget -q --show-progress https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.sh -O- | bash"
	} else if _, err := exec.LookPath("curl"); err == nil {
		installCmd = "curl -L https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.sh | bash"
	} else {
		cli.Error("Please install either wget or curl.")
		os.Exit(1)
	}

	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		cli.Error("Failed to update codegame-cli.")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
