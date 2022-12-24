package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/code-game-project/go-utils/exec"
	"github.com/code-game-project/go-utils/external"
	"github.com/spf13/cobra"
)

var cgLSPCGEPath = filepath.Join(xdg.DataHome, "codegame", "bin", "lsp", "cge")

// cgeCmd represents the cge command
var cgeCmd = &cobra.Command{
	Use:   "cge",
	Short: "Launch cge-ls.",
	Args:  cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
		version, err := external.LatestGithubTag("code-game-project", "cg-gen-events")
		abort(err)
		version = strings.TrimPrefix(version, "v")

		exeName, err := external.InstallProgram("cge-ls", "cge-ls", "https://github.com/code-game-project/cg-gen-events", version, cgLSPCGEPath)
		abort(err)

		_, err = exec.Execute(false, filepath.Join(cgLSPCGEPath, exeName), args...)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	lspCmd.AddCommand(cgeCmd)
}
