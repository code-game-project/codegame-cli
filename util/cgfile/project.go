package cgfile

import (
	"os"
	"path/filepath"

	"github.com/Bananenpro/cli"
)

func FindProjectRoot() (string, error) {
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
			return "", cli.Error("Not in a CodeGame project directory")
		}
		dir = parent
	}
}

func FindProjectRootRelative() (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}
	wd, _ := os.Getwd()
	return filepath.Rel(wd, root)
}
