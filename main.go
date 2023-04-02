package main

import "github.com/code-game-project/codegame-cli/cmd"

// populated by CI (e.g. "1.2.3")
var version = "dev"

func main() {
	cmd.Execute(version)
}
