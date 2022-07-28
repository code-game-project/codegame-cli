package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Bananenpro/cli"
	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/pkg/cgfile"
	"github.com/code-game-project/codegame-cli/pkg/external"
	"github.com/code-game-project/codegame-cli/pkg/semver"
)

var modulesPath = filepath.Join(xdg.DataHome, "codegame", "bin", "codegame-cli", "modules")

func ExecuteNewClient(data NewClientData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return execute(data.Lang, data.LibraryVersion, "client", jsonData, "new", "client")
}

func ExecuteNewServer(data NewServerData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return execute(data.Lang, data.LibraryVersion, "server", jsonData, "new", "server")
}

func ExecuteUpdate(data UpdateData, cgData *cgfile.CodeGameFileData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return execute(cgData.Lang, data.LibraryVersion, cgData.Type, jsonData, "update")
}

func ExecuteRun(data RunData, cgData *cgfile.CodeGameFileData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return execute(cgData.Lang, "latest", cgData.Type, jsonData, "run")
}

func ExecuteBuild(data BuildData, cgData *cgfile.CodeGameFileData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return execute(cgData.Lang, "latest", cgData.Type, jsonData, "build")
}

func execute(lang, libraryVersion, projectType string, data []byte, args ...string) error {
	version, err := findModuleVersion(lang, libraryVersion, projectType)
	if err != nil {
		return err
	}

	exeName, err := installModule(lang, version)
	if err != nil {
		return err
	}

	programName := filepath.Join(modulesPath, lang, exeName)
	if _, err := exec.LookPath(programName); err != nil {
		return fmt.Errorf("'%s' ist not installed!", programName)
	}

	configFile, err := os.CreateTemp("", "codegame-cli-module-config-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(configFile.Name())

	_, err = configFile.Write(data)
	configFile.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command(programName, args...)
	cmd.Env = append(os.Environ(), "CONFIG_FILE="+configFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findModuleVersion(name, libraryVersion, projectType string) (string, error) {
	if libraryVersion == "latest" {
		version, err := external.LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	body, err := external.LoadVersionsJSON("code-game-project", "codegame-cli-"+name)
	if err != nil {
		cli.Warn("Couldn't fetch versions.json. Using latest %s module version.", name)
		version, err := external.LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}
	defer body.Close()

	type jsonObj struct {
		Server map[string]string
		Client map[string]string
	}

	var versions jsonObj
	err = json.NewDecoder(body).Decode(&versions)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest %s module version.", name)
		version, err := external.LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	var v string
	if projectType == "client" {
		v = semver.CompatibleVersion(versions.Client, libraryVersion)
	} else if projectType == "server" {
		v = semver.CompatibleVersion(versions.Server, libraryVersion)
	} else {
		return "", errors.New("invalid project type")
	}

	if v == "latest" {
		version, err := external.LatestGithubTag("code-game-project", "codegame-cli-"+name)
		return strings.TrimPrefix(version, "v"), err
	}

	v, err = external.GithubTagFromVersion("code-game-project", "codegame-cli-"+name, v)
	return strings.TrimPrefix(v, "v"), err
}

func installModule(name, version string) (string, error) {
	path := filepath.Join(modulesPath, name)
	return external.InstallProgram(fmt.Sprintf("codegame-cli-%s", name), fmt.Sprintf("codegame-%s", name), fmt.Sprintf("https://github.com/code-game-project/codegame-cli-%s", name), version, path)
}
