package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/external"
	"github.com/code-game-project/go-utils/server"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display some information about a game server.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) > 0 {
			url = strings.ToLower(args[0])
		} else {
			var err error
			url, err = cli.Input("Game server URL:")
			abort(err)
		}
		api, err := server.NewAPI(url)
		abort(err)

		info, err := api.FetchGameInfo()
		abort(err)

		printInfo(info)

		cli.PrintColor(cli.Yellow, "\nTo view the documentation of this game run:\n%s docs %s", os.Args[0], external.TrimURL(url))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func printInfo(info server.GameInfo) {
	out := colorable.NewColorableStdout()
	printInfoProperty(out, "Display Name", info.DisplayName, 17)
	printInfoProperty(out, "Name", info.Name, 17)
	printInfoProperty(out, "Description", info.Description, 17)
	printInfoProperty(out, "Version", info.Version, 17)
	printInfoProperty(out, "CodeGame Version", info.CGVersion, 17)
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
