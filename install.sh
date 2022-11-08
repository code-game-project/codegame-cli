#!/bin/bash

echo "Installing codegame-cli..."

cd /tmp

rm -f "codegame-cli.tar.gz"

os=$(uname)
arch=$(uname -m)

download () {
	if hash wget 2>/dev/null; then
		wget -q --show-progress https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-$1-$2.tar.gz -O codegame-cli.tar.gz || exit 1
	elif hash curl 2>/dev/null; then
		curl -L https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-$1-$2.tar.gz > codegame-cli.tar.gz || exit 1
	else
		echo "Please install either wget or curl."
		exit 1
	fi
}

shopt -s nocasematch

if [[ $os == *"linux"* ]]; then
	if [[ $arch == *"x86"* ]]; then
		echo "Detected OS: Linux x86_64"
		download "linux" "amd64"
	elif [[ $arch == *"aarch64"* ]]; then
		echo "Detected OS: Linux ARM64"
		download "linux" "arm64"
	elif [[ $arch == *"arm"* ]]; then
		echo "Detected OS: Linux ARM64"
		download "linux" "arm64"
	else
		echo "Detected OS: $os $arch"
		echo "Your architecture is not supported by this installer."
		exit 1
	fi
elif [[ $os == *"darwin"* ]]; then
	if [[ $arch == *"x86"* ]]; then
		echo "Detected OS: macOS x86_64"
		download "darwin" "amd64"
	elif [[ $arch == *"aarch64"* ]]; then
		echo "Detected OS: macOS ARM64"
		download "darwin" "arm64"
	elif [[ $arch == *"arm"* ]]; then
		echo "Detected OS: macOS ARM64"
		download "darwin" "arm64"
	else
		echo "Detected OS: $os $arch"
		echo "Your architecture is not supported by this installer."
		exit 1
	fi
else
	echo "Detected OS: $os $arch"
	echo "Your OS is not supported by this installer."
	exit 1
fi

if [[ :$PATH: == *:"$HOME/.local/bin":* ]] ; then
	echo "Installing binaries into ~/.local/bin..."
	mkdir -p $HOME/.local/bin || exit 1
	if test -f /usr/local/bin/codegame; then
		echo "Removing old version in /usr/local/bin..."
		sudo rm -f /usr/local/bin/codegame
	fi
	tar -xzf codegame-cli.tar.gz codegame && mv codegame $HOME/.local/bin || exit 1
else
	echo "Installing binaries into /usr/local/bin..."
	sudo mkdir -p /usr/local/bin || exit 1
	if test -f $HOME/.local/bin/codegame; then
		echo "Removing old version in ~/.local/bin..."
		rm -f $HOME/.local/bin/codegame
	fi
	tar -xzf codegame-cli.tar.gz codegame && sudo mv codegame /usr/local/bin || exit 1
fi

rm codegame-cli.tar.gz

echo "Done."
