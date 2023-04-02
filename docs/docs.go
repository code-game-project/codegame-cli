package docs

import (
	"fmt"

	"github.com/code-game-project/cli-utils/exec"
)

const DocsURL = "https://code-game.org/docs/intro"

func Open() error {
	fmt.Printf("Opening %s in default web browser...\n", DocsURL)
	err := exec.OpenInBrowser(DocsURL)
	if err != nil {
		return fmt.Errorf("failed to open documentation: %w", err)
	}
	return nil
}
