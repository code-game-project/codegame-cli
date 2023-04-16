package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/feedback"
	"github.com/code-game-project/cli-utils/request"
	"github.com/code-game-project/cli-utils/server"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about a game server",
	Run: func(_ *cobra.Command, args []string) {
		var url string
		if len(args) == 0 {
			file, err := cgfile.Load("")
			if err == nil {
				url = file.URL
			}
			url = cli.Input("Game URL:", true, url)
		} else {
			url = args[0]
		}
		url = request.TrimURL(url)
		info, err := server.FetchGameInfo(url)
		if err != nil {
			feedback.Fatal("codegame-cli", "%s is not a valid CodeGame game server or is not reachable", url)
			os.Exit(1)
		}
		printInfo(url, info)
	},
}

func printInfo(url string, info server.GameInfo) {
	out := colorable.NewColorableStdout()
	printInfoProperty(out, "Game URL", request.TrimURL(url), 17)
	printInfoProperty(out, "Display Name", info.DisplayName, 17)
	printInfoProperty(out, "Name", info.Name, 17)
	printInfoProperty(out, "Description", info.Description, 17)
	printInfoProperty(out, "Version", info.Version, 17)
	printInfoProperty(out, "CodeGame Version", info.CGVersion.String(), 17)
	printInfoProperty(out, "Repository", info.RepositoryURL, 17)
}

func printInfoProperty(out io.Writer, name, value string, labelWidth int) {
	if value == "" {
		return
	}

	label := name + ":"
	if labelWidth-utf8.RuneCountInString(label) > 0 {
		label += strings.Repeat(" ", labelWidth-utf8.RuneCountInString(label))
	}

	fmt.Fprintf(out, "\x1b[36m%s\x1b[0m %s\n", label, value)
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
