package cmd

import (
	"os"
	"path/filepath"

	_ "embed"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/go-utils/cggenevents"
	"github.com/code-game-project/go-utils/exec"
	"github.com/code-game-project/go-utils/server"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
)

//go:embed templates/css/docs.css
var docsStyle string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "View the documention of CodeGame or a specific game in your webbrowser.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Print("Opening documentation...")
			err := exec.OpenBrowser("https://docs.code-game.org")
			abort(err)
		}

		api, err := server.NewAPI(args[0])
		abort(err)

		url := api.BaseURL()

		cge, err := api.GetCGEFile()
		abort(err)

		cgeVersion, err := cggenevents.ParseCGEVersion(cge)
		abort(err)

		err = cggenevents.CGGenEvents(cgeVersion, os.TempDir(), url, "markdown")
		abort(err)

		md, err := os.ReadFile(filepath.Join(os.TempDir(), "event_docs.md"))
		abort(err)

		md = markdown.NormalizeNewlines(md)
		text := markdown.ToHTML(md, parser.NewWithExtensions(parser.CommonExtensions|parser.AutoHeadingIDs), html.NewRenderer(html.RendererOptions{
			CSS:   filepath.Join(os.TempDir(), "event_docs.css"),
			Flags: html.CommonFlags | html.CompletePage,
		}))

		os.Remove(filepath.Join(os.TempDir(), "event_docs.md"))

		err = os.WriteFile(filepath.Join(os.TempDir(), "event_docs.html"), text, 0644)
		abort(err)

		err = os.WriteFile(filepath.Join(os.TempDir(), "event_docs.css"), []byte(docsStyle), 0644)
		abort(err)

		cli.Print("Opening documentation...")

		err = exec.OpenBrowser(filepath.Join(os.TempDir(), "event_docs.html"))
		abort(err)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
