package external

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
)

var cgGenEventsPath = filepath.Join(xdg.DataHome, "codegame", "bin", "cg-gen-events")

func InstallCGGenEvents(cgeVersion string) error {
	exeName := fmt.Sprintf("cg-gen-events-%s", strings.ReplaceAll(cgeVersion, ".", "_"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cgGenEventsPath, exeName)); err == nil {
		return nil
	}

	version, err := GithubTagFromVersion("code-game-project", "cg-gen-events", cgeVersion)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("cg-gen-events-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		filename = fmt.Sprintf("cg-gen-events-%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	}

	res, err := http.Get(fmt.Sprintf("https://github.com/code-game-project/cg-gen-events/releases/download/%s/%s", version, filename))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = os.MkdirAll(cgGenEventsPath, 0755)
	if err != nil {
		return err
	}

	return UntargzFile(res.Body, "cg-gen-events", filepath.Join(cgGenEventsPath, exeName))
}

func CGGenEvents(outputDir, url, cgeVersion, language string) error {
	exeName := fmt.Sprintf("cg-gen-events-%s", strings.ReplaceAll(cgeVersion, ".", "_"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cgGenEventsPath, exeName)); err != nil {
		err = InstallCGGenEvents(cgeVersion)
		if err != nil {
			return err
		}
	}

	return Execute(filepath.Join(cgGenEventsPath, exeName), url, "--languages", language, "--output", outputDir)
}
