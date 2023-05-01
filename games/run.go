package games

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/exec"
	"github.com/code-game-project/cli-utils/feedback"
	"github.com/code-game-project/cli-utils/modules"
)

var gamesDir = filepath.Join(xdg.DataHome, "codegame", "games")

func init() {
	os.MkdirAll(gamesDir, 0o755)
}

func Run(repoURL string, args []string) error {
	path := filepath.Join(gamesDir, url.PathEscape(repoURL))
	if _, err := os.Stat(path); err == nil {
		return run(path, args)
	}
	err := exec.ExecuteDimmed("git", "clone", repoURL, path)
	if err != nil {
		return fmt.Errorf("clone game repository: %w", err)
	}
	return run(path, args)
}

func run(path string, args []string) error {
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("chdir to game repository: %w", err)
	}
	file, err := cgfile.Load("")
	if err != nil {
		os.RemoveAll(path)
		return fmt.Errorf("not a CodeGame project")
	}
	if file.ProjectType != "server" {
		os.RemoveAll(path)
		return fmt.Errorf("not a game server")
	}
	mod, err := modules.LoadModule(file.Language)
	if err != nil {
		return fmt.Errorf("load module: %w", err)
	}
	return mod.ExecRunServer(file.ModVersion, mod.Lang, nil, args)
}

func ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(gamesDir)
	if err != nil {
		return nil, fmt.Errorf("read games dir: %w", err)
	}

	games := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			repo, err := url.PathUnescape(e.Name())
			if err != nil {
				feedback.Warn("codegame-cli", "Invalid game repository URL %s: %w", e.Name(), err)
				continue
			}
			games = append(games, repo)
		}
	}
	return games, nil
}
