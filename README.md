# codegame-cli
![CG Server Version](https://img.shields.io/badge/GameServer-v0.1+-yellow)
![CG Client Version](https://img.shields.io/badge/Client-v0.3+-yellow)

The official [CodeGame](https://code-game.org) CLI.

## Usage

View the CodeGame documentation:
```sh
codegame docs
```

Create a new project:
```sh
codegame new
```

Run a project:
```sh
codegame run
```

Get information about a game server:
```sh
codegame info <url>
```

View the documentation of a game:
```sh
codegame docs <url>
```

Help:
```sh
codegame --help
```

## Features

- View the CodeGame documentation
- View information about a game server
- Automatic project setup
  - CodeGame version detection and management
  - Create a new client
    - Go
  - Create a new game server
    - Go
  - Initialize Git
  - Create README and LICENSE files

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
