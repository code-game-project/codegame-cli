package commands

import (
	"os"
	"path/filepath"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/util"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/ogier/pflag"

	_ "embed"
)

//go:embed templates/css/docs.css
var docsStyle string

func Docs() error {
	if pflag.NArg() == 1 {
		cli.Begin("Opening documentation...")
		err := util.OpenBrowser("https://docs.code-game.org")
		if err != nil {
			cli.Error(err.Error())
		}
		cli.Finish()
		return err
	}

	cli.Begin("Generating markdown documentation...")

	url := baseURL(pflag.Arg(1))

	cgeVersion, err := util.GetCGEVersion(url)
	if err != nil {
		return err
	}

	err = util.CGGenEvents(os.TempDir(), url, cgeVersion, "markdown")
	if err != nil {
		return err
	}

	cli.Finish()

	cli.Begin("Converting documentation to HTML...")

	md, err := os.ReadFile(filepath.Join(os.TempDir(), "event_docs.md"))
	if err != nil {
		return cli.Error(err.Error())
	}

	md = markdown.NormalizeNewlines(md)
	text := markdown.ToHTML(md, parser.NewWithExtensions(parser.CommonExtensions|parser.AutoHeadingIDs), html.NewRenderer(html.RendererOptions{
		CSS:   filepath.Join(os.TempDir(), "event_docs.css"),
		Flags: html.CommonFlags | html.CompletePage,
	}))

	os.Remove(filepath.Join(os.TempDir(), "event_docs.md"))

	err = os.WriteFile(filepath.Join(os.TempDir(), "event_docs.html"), text, 0644)
	if err != nil {
		return cli.Error(err.Error())
	}

	err = os.WriteFile(filepath.Join(os.TempDir(), "event_docs.css"), []byte(docsStyle), 0644)
	if err != nil {
		return cli.Error(err.Error())
	}

	cli.Finish()

	cli.Begin("Opening documentation...")

	err = util.OpenBrowser(filepath.Join(os.TempDir(), "event_docs.html"))
	if err != nil {
		cli.Error(err.Error())
	}

	cli.Finish()

	return nil
}