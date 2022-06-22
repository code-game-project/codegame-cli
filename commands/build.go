package commands

import (
	"os"

	"github.com/Bananenpro/cli"
	"github.com/code-game-project/codegame-cli/util/cgfile"
	"github.com/code-game-project/codegame-cli/util/modules"
)

func Build() error {
	rootRelative, err := cgfile.FindProjectRootRelative()
	if err != nil {
		return err
	}

	data, err := cgfile.LoadCodeGameFile(rootRelative)
	if err != nil {
		return cli.Error("Failed to load .codegame.json")
	}

	args := []string{"build"}
	if len(os.Args) > 2 {
		args = append(args, os.Args[2:]...)
	}

	switch data.Lang {
	case "go":
		err = modules.Execute("go", "latest", data.Type, args...)
	default:
		return cli.Error("'build' is not supported for '%s'", data.Lang)
	}
	if err != nil {
		return err
	}
	return nil
}
