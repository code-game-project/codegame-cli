# codegame-cli
![CG Server Version](https://img.shields.io/badge/GameServer-v0.1+-yellow)
![CG Client Version](https://img.shields.io/badge/Client-v0.3+-yellow)

The official [CodeGame](https://code-game.org) CLI.

## Commands

### Docs and info

View the CodeGame documentation:
```sh
codegame docs
```

View the documentation of a game:
```sh
codegame docs <url>
```

Get information about a game server:
```sh
codegame info <url>
```

### Project templates

Create a new project:
```sh
codegame new
```

Update event definitions, wrappers and libraries to match the latest game version:
```sh
codegame update
```

Permanently switch to a different game URL:
```sh
codegame change-url <new_url>
```

### Running and building

Run a project:
```sh
codegame run
```

Build a project:
```sh
codegame build
```

### Help

Display help:
```sh
codegame --help
```

## Installation

### Windows

1. Open the Start menu
2. Search for `powershell`
3. Hit `Run as Administrator`
4. Paste the following commands and hit enter:

#### Install

```powershell
Invoke-WebRequest -Uri "https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-windows-amd64.zip" -OutFile "C:\Program Files\codegame-cli.zip"
Expand-Archive -LiteralPath "C:\Program Files\codegame-cli.zip" -DestinationPath "C:\Program Files\codegame-cli"
rm "C:\Program Files\codegame-cli.zip"
Set-ItemProperty -Path 'Registry::HKEY_LOCAL_MACHINE\System\CurrentControlSet\Control\Session Manager\Environment' -Name PATH -Value "$((Get-ItemProperty -Path 'Registry::HKEY_LOCAL_MACHINE\System\CurrentControlSet\Control\Session Manager\Environment' -Name PATH).path);C:\Program Files\codegame-cli"
```

**IMPORTANT:** Please reboot for the installation to take effect.

#### Update

```powershell
rm -r -fo "C:\Program Files\codegame-cli"
Invoke-WebRequest -Uri "https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-windows-amd64.zip" -OutFile "C:\Program Files\codegame-cli.zip"
Expand-Archive -LiteralPath "C:\Program Files\codegame-cli.zip" -DestinationPath "C:\Program Files\codegame-cli"
rm "C:\Program Files\codegame-cli.zip"
```

### macOS

Open the Terminal application, paste the command for your architecture and hit enter.

To update, simply run the command again.

#### x86_64

```sh
curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-darwin-amd64.tar.gz | tar -xz codegame && sudo mv codegame /usr/local/bin
```

#### ARM64

```sh
curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-darwin-arm64.tar.gz | tar -xz codegame && sudo mv codegame /usr/local/bin
```

### Linux

Open a terminal, paste the command for your architecture and hit enter.

To update, simply run the command again.

#### x86_64

```sh
curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-linux-amd64.tar.gz | tar -xz codegame && sudo mv codegame /usr/local/bin
```

#### ARM64

```sh
curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-linux-arm64.tar.gz | tar -xz codegame && sudo mv codegame /usr/local/bin
```

### Other

You can download a prebuilt binary file for your operating system on the [releases](https://github.com/code-game-project/codegame-cli/releases) page.

### Compiling from source

#### Prerequisites

- [Go](https://go.dev/) 1.18+

```sh
git clone https://github.com/code-game-project/codegame-cli.git
cd codegame-cli
go build .
```

## Modules

Language specific functionality is provided by modules.
*codegame-cli* automatically downloads the correct module version corresponding to the CodeGame version and
executes the downloaded binary at the root of the project directory with the respective command.

Additionally, the _CONFIG_FILE_ environment variable is set with the path to a temporary file containing command specific configuration in JSON format.

Modules can be implemented in any language. However, it is recommended to write them in Go
using the `github.com/code-game-project/codegame-cli/pkg/*` packages and the `github.com/Bananenpro/cli` package for CLI interaction in order to be consistent with the CLI and other modules.

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

Runs the project with the specified command line arguments.

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

Copyright (c) 2022 CodeGame Contributors (https://code-game.org/contributors)

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
