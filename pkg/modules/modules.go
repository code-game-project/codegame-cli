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
	"github.com/code-game-project/codegame-cli/pkg/external"
	"github.com/code-game-project/codegame-cli/pkg/semver"
)

var modulesPath = filepath.Join(xdg.DataHome, "codegame", "bin", "codegame-cli", "modules")

func Execute(name, libraryVersion, projectType string, args ...string) error {
	version, err := findModuleVersion(name, libraryVersion, projectType)
	if err != nil {
		return err
	}

	exeName, err := installModule(name, version)
	if err != nil {
		return err
	}

	programName := filepath.Join(modulesPath, name, exeName)
	if _, err := exec.LookPath(programName); err != nil {
		return fmt.Errorf("'%s' ist not installed!", programName)
	}
	cmd := exec.Command(programName, args...)
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
