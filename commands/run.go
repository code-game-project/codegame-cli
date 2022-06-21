package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/Bananenpro/pflag"
	"github.com/code-game-project/codegame-cli/util/cgfile"
)

func Run() error {
	flagSet := pflag.NewFlagSet("run", pflag.ExitOnError)
	flagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags:           true,
		PassUnknownFlagsToArgs: true,
	}

	var overrideURL string
	flagSet.StringVarP(&overrideURL, "override-url", "u", "", "Override the game URL in .codegame.json for this specific run.")
	var help bool
	flagSet.BoolVarP(&help, "help", "h", false, "Show help.")

	if help {
		fmt.Fprintf(os.Stderr, "Usage: %s run [...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flagSet.PrintDefaults()
		os.Exit(1)
	}
	flagSet.Parse(os.Args[2:])

	root, err := findProjectRoot()
	if err != nil {
		return cli.Error("Not in a CodeGame project directory")
	}

	data, err := cgfile.LoadCodeGameFile(root)
	if err != nil {
		return cli.Error("Failed to load .codegame.json")
	}

	url := data.URL
	if overrideURL != "" {
		url = overrideURL
	}

	var cmdName string
	args := []string{}
	switch data.Lang {
	case "go":
		cmdName = "go"
		args = append(args, "run")
	default:
		return cli.Error("'run' is not supported for '%s'", data.Lang)
	}
	args = append(args, flagSet.Args()...)

	env := make([]string, 0, len(os.Environ())+1)
	env = append(env, "CG_GAME_URL="+url)
	env = append(env, os.Environ()...)

	if _, err := exec.LookPath(cmdName); err != nil {
		cli.Error("'%s' ist not installed!", cmdName)
		return err
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	err = cmd.Run()
	if err != nil {
		cli.Error("Failed to run 'CG_GAME_URL=%s %s %s'", url, cmdName, strings.Join(args, " "))
	}
	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return "", err
		}
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() == ".codegame.json" {
				return dir, nil
			}
		}

		parent := filepath.Dir(filepath.Clean(dir))
		if parent == dir {
			return "", errors.New("not found")
		}
		dir = parent
	}
}
