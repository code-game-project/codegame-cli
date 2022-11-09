package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall codegame-cli.",
	Run: func(cmd *cobra.Command, args []string) {
		shallow, err := cmd.Flags().GetBool("shallow")
		abort(err)
		removeData, err := cmd.Flags().GetBool("remove-data")
		abort(err)

		yes, err := cli.YesNo("Are you sure you want to uninstall codegame-cli?", false)
		abort(err)
		if !yes {
			cli.Print("Canceled.")
			return
		}

		dataDir := filepath.Join(xdg.DataHome, "codegame")
		if removeData {
			cli.Print("Removing session files...")
			err := os.RemoveAll(filepath.Join(dataDir, "games"))
			abortf("Failed to remove session files: %s", err)
		}
		if !shallow {
			cli.Print("Removing installed tools...")
			err := os.RemoveAll(filepath.Join(dataDir, "bin"))
			abortf("Failed to remove tools: %s", err)
		}
		if removeData && !shallow {
			os.RemoveAll(dataDir)
		}

		cli.Print("Removing codegame-cli...")
		switch runtime.GOOS {
		case "windows":
			uninstallCLIWindows()
		default:
			uninstallCLIUnix()
		}

		cli.Print("Removing cache files...")
		os.RemoveAll(filepath.Join(xdg.CacheHome, "codegame", "cli"))
		if files, err := os.ReadDir(filepath.Join(xdg.CacheHome, "codegame")); err == nil && len(files) == 0 {
			os.RemoveAll(filepath.Join(xdg.CacheHome, "codegame"))
		}

		cli.Success("Successfully removed codegame-cli from your system!")
	},
}

func uninstallCLIWindows() {
	homeDir, err := os.UserHomeDir()
	abortf("Failed to get the user home directory: %s", err)
	homeDir = filepath.Clean(homeDir)

	installDir := homeDir + "\\AppData\\Local\\Programs\\codegame-cli"

	cli.BeginLoading("Removing codegame-cli from PATH...")
	cmd := exec.Command("Powershell.exe", "-Command", "[System.Environment]::SetEnvironmentVariable(\"PATH\", [System.Environment]::GetEnvironmentVariable(\"PATH\",\"USER\") -replace \";"+strings.ReplaceAll(installDir, "\\", "\\\\")+"\",\"USER\")")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	abortf("Failed to remove codegame-cli from PATH: %s", err)
	cli.FinishLoading()

	err = os.RemoveAll(installDir)
	abortf("Failed to uninstall codegame-cli: %s", err)
}

func uninstallCLIUnix() {
	homeDir, err := os.UserHomeDir()
	abortf("Failed to get the user home directory: %s", err)
	homeDir = filepath.Clean(homeDir)

	if _, err := os.Stat("/usr/local/bin/codegame"); !os.IsNotExist(err) {
		cmd := exec.Command("bash", "-c", "sudo rm /usr/local/bin/codegame")
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		abortf("Failed to uninstall codegame-cli: %s", err)
	}
	if _, err := os.Stat(homeDir + "/.local/bin/codegame"); !os.IsNotExist(err) {
		err := os.Remove(homeDir + "/.local/bin/codegame")
		abortf("Failed to uninstall codegame-cli: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolP("shallow", "s", false, "Don't remove CodeGame tools installed by codegame-cli, e.g. cg-gen-events, cg-debug.")
	uninstallCmd.Flags().BoolP("remove-data", "d", false, "Remove game session files.")
}
