package main

import (
	"github.com/code-game-project/codegame-cli/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
