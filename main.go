package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/code-game-project/codegame-cli/cli"
	"github.com/code-game-project/codegame-cli/commands"
	"github.com/code-game-project/codegame-cli/external"
	"github.com/ogier/pflag"
)

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nDescription:")
		fmt.Fprintln(os.Stderr, "\nThe official CodeGame CLI.")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "\tnew \tCreate a new project.")
		fmt.Fprintln(os.Stderr, "\tupdate \tUpdate the CodeGame libraries in the current project.")
		fmt.Fprintln(os.Stderr, "\tdocs \tOpen the CodeGame documentation in a web browser.")
		fmt.Fprintln(os.Stderr, "\nAbout: https://github.com/code-game-project")
		fmt.Fprintln(os.Stderr, "Copyright (c) 2022 CodeGame Contributors (https://github.com/orgs/code-game-project/people)")
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if pflag.NArg() == 0 {
		pflag.Usage()
		os.Exit(1)
	}

	command := strings.ToLower(pflag.Arg(0))

	var err error
	switch command {
	case "new":
		err = commands.New()
		if err == cli.ErrCanceled {
			cli.Info(err.Error())
		}
	case "docs":
		err = external.OpenBrowser("https://github.com/code-game-project/docs/blob/main/README.md")
		if err != nil {
			cli.Error(err.Error())
		}
	default:
		cli.Error("Unknown command: %s", strings.ToLower(pflag.Arg(0)))
		pflag.Usage()
		os.Exit(1)
	}
	if err != nil {
		if cli.IsInProgress() {
			fmt.Println()
		}
		os.Exit(1)
	}
}
