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
		fmt.Fprintln(os.Stderr, "The official CodeGame CLI.")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "\tnew \tCreate a new project.")
		fmt.Fprintln(os.Stderr, "\tdocs \tOpen the CodeGame documentation in a web browser.")
		fmt.Fprintln(os.Stderr, "\nAbout: https://code-game.org")
		fmt.Fprintln(os.Stderr, "Copyright (c) 2022 CodeGame Contributors (https://code-game.org/contributors)")
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
	case "info":
		err = commands.Info()
	case "docs":
		err = external.OpenBrowser("https://docs.code-game.org")
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
