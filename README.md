# codegame-cli
![CodeGame Version](https://img.shields.io/badge/CodeGame-v0.6+-orange)

The official [CodeGame](https://code-game.org) CLI.

**[Install](#installation)**

## Commands

### Docs and info

View the CodeGame documentation:
```
codegame docs
```

View the documentation of a game:
```
codegame docs <url>
```

Get information about a game server:
```
codegame info <url>
```

### Project templates

Create a new project:
```
codegame new
```

Update event definitions, wrappers and libraries to match the latest game version:
```
codegame update
```

Permanently switch to a different game URL:
```
codegame change-url <new_url>
```

### Running and building

Run a project:
```
codegame run
```

Build a project:
```
codegame build
```

### Session management

List all sessions:
```
codegame session list
```

Show session details:
```
codegame session show
```

Remove a session:
```
codegame session remove
```

Export a session to [CodeGame Share](https://share.code-game.org):
```
codegame session export
```

Import a session from [CodeGame Share](https://share.code-game.org):
```
codegame session import <id>
```

### Games

List all games on a server:
```
codegame game list <url>
```

Create a new game on a server:
```
codegame game create <url>
```

### Sharing with [CodeGame Share](https://share.code-game.org)

Share a game:
```
codegame share game
```

Share a spectate link:
```
codegame share spectate
```

Share a session:
```
codegame share session
```

### cg-gen-events

Download and execute the correct version of [cg-gen-events](https://github.com/code-game-project/cg-gen-events):
```
codegame gen-events <input>
```

### cg-debug

Download and execute the correct version of [cg-debug](https://github.com/code-game-project/cg-debug):
```
codegame debug <url>
```

### Completion

Generate an autocompletion script for codegame-cli for the specified shell:
```
codegame completion <bash|zsh|fish|powershell>
```

### Help

Display general help:
```
codegame --help
```

Display help about a specific command:
```
codegame help <cmd>
```

## Installation

### Windows

[Download](https://github.com/code-game-project/codegame-cli/releases/latest/download/install.bat) and execute `install.bat`.

### macOS/Linux

Paste one of the following commands into a terminal window:

#### curl

```bash
curl -L https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.sh | bash
```

#### wget (in case curl is not installed)

```bash
wget -q --show-progress https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.sh -O- | bash
```

## Uninstallation

To remove codegame-cli from your system run:

```bash
codegame uninstall
```

### Compiling from source

#### Prerequisites

- [Go](https://go.dev/) 1.18+

```
git clone https://github.com/code-game-project/codegame-cli.git
cd codegame-cli
go build -o codegame .
```

## Modules

Language specific functionality is provided by modules.
*codegame-cli* automatically downloads the correct module version corresponding to the CodeGame version and
executes the downloaded binary at the root of the project directory with the respective command.

Additionally, the _CONFIG_FILE_ environment variable is set with the path to a temporary file containing command specific configuration in JSON format.

Modules can be implemented in any language. However, it is recommended to write them in Go
using the `github.com/code-game-project/go-utils` and `github.com/Bananenpro/cli` packages in order to avoid bugs and be consistent with the CLI and other modules.

### new

Creates a new project of the type provided by a second command line argument:

#### client

Creates a new game client.
This includes integration with the client library for the language, a functioning template with a main file and language specific metadata files like `package.json` or similar
and wrappers around the library to make its usage easier including setting the game URL with the CG_GAME_URL environment variable (required).

##### config data

```jsonc
{
	"lang": "go", // the chosen programming language (in case one module supports multiple languages)
	"name": "my_game", // the name of the game
	"url": "my_game.example.com", // the URL of the game server
	"library_version": "0.9.2" // the version of the client library to use
}
```

#### server

Creates a new game server.
This includes integration with the server library of the language and a functioning template, which implements common logic like starting the server and providing a game class.

##### config data

```jsonc
{
	"lang": "go", // the chosen programming language (in case one module supports multiple languages)
	"library_version": "0.9.2" // the version of the server library to use
}
```

### update

Updates the library to the specified version and all other dependencies used by the project to their newest compatible version.
Additionally all wrappers are updated.

##### config data

```jsonc
{
	"lang": "go", // the chosen programming language (in case one module supports multiple languages)
	"library_version": "0.9.2" // the new version of the library to use
}

```

### run

Runs the project with the specified command line arguments and the *CG_GAME_URL* environment variable set to the URL specified in the `.codegame.json` file.

##### config data

```jsonc
{
	"lang": "go", // the chosen programming language (in case one module supports multiple languages)
	"args": ["-f", "my_file.txt"] // the command line arguments for the application
}
```

### build

Builds the projects and injects the URL specified in the `.codegame.json` file, which makes the *CG_GAME_URL* environment variable optional.

##### config data

```jsonc
{
	"lang": "go", // the chosen programming language (in case one module supports multiple languages)
	"output": "bin/my-game" // The name of the output file
}
```

## License

Copyright (c) 2022 Julian Hofmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
