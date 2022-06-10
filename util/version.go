package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/code-game-project/codegame-cli/cli"
)

const currentVersion = "0.5.2"

// CheckVersion prints a warning, if there is a newer version of codegame-cli available.
func CheckVersion() {
	latest, err := getLatestVersion()
	if err != nil {
		return
	}

	currentMajor, currentMinor, currentPatch, err := ParseVersion(currentVersion)
	if err != nil {
		return
	}

	latestMajor, latestMinor, latestPatch, err := ParseVersion(latest)
	if err != nil {
		return
	}

	if latestMajor > currentMajor || (latestMajor == currentMajor && latestMinor > currentMinor) || (latestMajor == currentMajor && latestMinor == currentMinor && latestPatch > currentPatch) {
		cli.Warn("You are using an old version of codegame-cli (v%s).\nUpdate to the latest version (v%s): https://github.com/code-game-project/codegame-cli#installation", currentVersion, latest)
	}
}

// CompatibleVersion returns the next best compatible version in the versions map.
func CompatibleVersion(versions map[string]string, version string) string {
	// check exact match
	if v, ok := versions[version]; ok {
		return v
	}

	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		cli.Warn("Invalid versions.json. Using latest version.")
		return "latest"
	}
	major := parts[0]

	// get all minor versions of the requested major version
	compatibleMinorVersions := make([]int, 0)
	for v := range versions {
		clientParts := strings.Split(v, ".")
		if len(clientParts) < 2 {
			cli.Warn("Invalid versions.json. Using latest version.")
			return "latest"
		}
		clientMajor := clientParts[0]
		if major == clientMajor {
			minor, err := strconv.Atoi(clientParts[1])
			if err != nil {
				cli.Warn("Invalid versions.json. Using latest version.")
				return "latest"
			}
			compatibleMinorVersions = append(compatibleMinorVersions, minor)
		}
	}
	if len(compatibleMinorVersions) == 0 {
		cli.Warn("No compatible version found. Using version.")
		return "latest"
	}

	minorStr := parts[1]
	minor, err := strconv.Atoi(minorStr)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest version.")
		return "latest"
	}

	// check closest minor version above requested
	closestMinor := -1
	for _, v := range compatibleMinorVersions {
		if v > minor && (closestMinor == -1 || closestMinor-minor > v-minor) {
			closestMinor = v
		}
	}
	if closestMinor >= 0 {
		v := versions[fmt.Sprintf("%s.%d", major, closestMinor)]
		cli.Warn("No exact version match found. Using version %s.", v)
		return v
	}

	// check closest minor version below requested
	closestMinor = -1
	for _, v := range compatibleMinorVersions {
		if v < minor && (closestMinor == -1 || minor-closestMinor > minor-v) {
			closestMinor = v
		}
	}
	if closestMinor >= 0 {
		v := versions[fmt.Sprintf("%s.%d", major, closestMinor)]
		cli.Warn("No exact version match found. Using version %s.", v)
		return v
	}

	cli.Warn("No compatible version found. Using latest version.")
	return "latest"
}

func getLatestVersion() (string, error) {
	cacheDir := filepath.Join(xdg.CacheHome, "codegame", "cli")
	os.MkdirAll(cacheDir, 0755)

	content, err := os.ReadFile(filepath.Join(cacheDir, "latest_version"))
	if err == nil {
		parts := strings.Split(string(content), "\n")
		if len(parts) >= 2 {
			cacheTime, err := strconv.Atoi(parts[0])
			if err == nil && time.Now().Unix()-int64(cacheTime) <= 3*24*60*60 {
				return parts[1], nil
			}
		}
	}

	tag, err := LatestGithubTag("code-game-project", "codegame-cli")
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(tag, "v")
	os.WriteFile(filepath.Join(cacheDir, "latest_version"), []byte(fmt.Sprintf("%d\n%s", time.Now().Unix(), version)), 0644)
	return version, nil
}

func ParseVersion(version string) (int, int, int, error) {
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")

	var major, minor, patch int
	var err error

	if len(parts) >= 1 {
		major, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version string: %s", version)
		}
	}

	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version string: %s", version)
		}
	}

	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version string: %s", version)
		}
	}

	return major, minor, patch, nil
}
