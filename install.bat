@echo off
echo Downloading latest installer...
Powershell.exe -Command "iwr -useb https://raw.githubusercontent.com/code-game-project/codegame-cli/main/install.ps1 | iex"
pause
