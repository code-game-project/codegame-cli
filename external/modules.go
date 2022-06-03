package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/cli"
)

var modulesPath = filepath.Join(xdg.DataHome, "codegame", "bin", "codegame-cli", "modules")

func ExecuteModule(dir, name, libraryVersion, projectType string, args ...string) error {
	version, err := findModuleVersion(name, libraryVersion, projectType)
	if err != nil {
		return cli.Error(err.Error())
	}

	exeName, err := installModule(name, version)
	if err != nil {
		return cli.Error(err.Error())
	}

	binaries, err := os.ReadDir(filepath.Join(modulesPath, name))
	if err != nil {
		return err
	}
	for _, b := range binaries {
		info, err := b.Info()
		if err == nil && info.Name() != exeName && strings.HasPrefix(info.Name(), fmt.Sprintf("codegame-%s_%s", name, strings.ReplaceAll(libraryVersion, ".", "-"))) {
			os.Remove(filepath.Join(cgGenEventsPath, info.Name()))
		}
	}

	return ExecuteInDir(dir, filepath.Join(modulesPath, name, exeName), args...)
}

func findModuleVersion(name, libraryVersion, projectType string) (string, error) {
	if libraryVersion == "latest" {
		version, err := LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	res, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/code-game-project/codegame-cli-%s/main/versions.json", name))
	if err != nil || res.StatusCode != http.StatusOK {
		cli.Warn("Couldn't fetch versions.json. Using latest go module version.")
		version, err := LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}
	defer res.Body.Close()

	type jsonObj struct {
		Server map[string]string
		Client map[string]string
	}

	var versions jsonObj
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest go module version.")
		version, err := LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	var version string
	if projectType == "client" {
		version = CompatibleVersion(versions.Client, libraryVersion)
	} else if projectType == "server" {
		version = CompatibleVersion(versions.Server, libraryVersion)
	} else {
		return "", errors.New("invalid project type")
	}

	if version == "latest" {
		version, err := LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	version, err = GithubTagFromVersion("code-game-project", "codegame-cli-"+name, version)
	return strings.TrimPrefix(version, "v"), err
}

func installModule(name, version string) (string, error) {
	exeName := fmt.Sprintf("codegame-%s_%s", name, strings.ReplaceAll(version, ".", "-"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	path := filepath.Join(modulesPath, name)

	if _, err := os.Stat(filepath.Join(path, exeName)); err == nil {
		return exeName, nil
	}

	filename := fmt.Sprintf("codegame-cli-%s-%s-%s.tar.gz", name, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		filename = fmt.Sprintf("codegame-cli-%s-%s-%s.zip", name, runtime.GOOS, runtime.GOARCH)
	}

	res, err := http.Get(fmt.Sprintf("https://github.com/code-game-project/codegame-cli-%s/releases/download/v%s/%s", name, version, filename))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		return exeName, UnzipFile(res.Body, "codegame-"+name+".exe", filepath.Join(path, exeName))
	}
	return exeName, UntargzFile(res.Body, "codegame-"+name, filepath.Join(path, exeName))
}
