if ((Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings').ProxyEnable) {
    $proxy = (Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings').ProxyServer
    $env:HTTP_PROXY = $proxy
    $env:HTTPS_PROXY = $proxy
	Write-Host "Using proxy: $proxy"
}

Write-Host "Installing codegame-cli..."

$InstallDir = "$HOME\AppData\Local\Programs\codegame-cli"
if (!($InstallDir | Test-Path)) {
	New-Item -ItemType "directory" -Path $InstallDir
	[System.Environment]::SetEnvironmentVariable("PATH",[System.Environment]::GetEnvironmentVariable("PATH","USER") + ";" + $InstallDir,"USER")

	Write-Host "Refreshing environment variables..."
	$HWND_BROADCAST = [intptr]0xffff;
	$WM_SETTINGCHANGE = 0x1a;
	$result = [uintptr]::zero
	if (-not ("win32.nativemethods" -As [type])) {
		Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
	}
	[void]([win32.nativemethods]::SendMessageTimeout($HWND_BROADCAST, $WM_SETTINGCHANGE, [uintptr]::Zero, "Environment", 2, 5000, [ref]$result))
	$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
} else {
	rm -r $InstallDir
	New-Item -ItemType "directory" -Path $InstallDir
}

$TempDir = [System.IO.Path]::GetTempPath()
cd $TempDir

Invoke-WebRequest -Uri https://github.com/code-game-project/codegame-cli/releases/latest/download/codegame-cli-windows-amd64.zip -OutFile codegame-cli.zip
Expand-Archive -LiteralPath codegame-cli.zip -DestinationPath $InstallDir
rm codegame-cli.zip


if (Get-Command code -ErrorAction SilentlyContinue) {
	Write-Host "Installing vscode-codegame..."
	Invoke-WebRequest -Uri https://github.com/code-game-project/vscode-codegame/releases/latest/download/codegame.vsix -OutFile .\codegame.vsix
	code --uninstall-extension code-game-project.codegame | Out-Null
	code --install-extension codegame.vsix
	rm codegame.vsix
}

Write-Host "Done." -ForegroundColor Green
