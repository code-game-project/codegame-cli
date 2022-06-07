package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/code-game-project/codegame-cli/cli"
)

func LatestGithubTag(owner, repo string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK || !HasContentType(res.Header, "application/json") {
		return "", fmt.Errorf("failed to access git tags from 'github.com/%s/%s'.", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", errors.New("failed to decode git tag data.")
	}
	return data[0].Name, nil
}

func GithubTagFromVersion(owner, repo, version string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK || !HasContentType(res.Header, "application/json") {
		return "", cli.Error("Couldn't access git tags from 'github.com/%s/%s'.", owner, repo)
	}
	defer res.Body.Close()
	type response []struct {
		Name string `json:"name"`
	}
	var data response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", cli.Error("Couldn't decode git tag data.")
	}

	for _, tag := range data {
		if strings.HasPrefix(tag.Name, "v"+version) {
			return tag.Name, nil
		}
	}
	return "", ErrTagNotFound
}

func LibraryVersionFromCGVersion(owner, repo, cgVersion string) string {
	res, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/versions.json", owner, repo))
	if err != nil || res.StatusCode != http.StatusOK {
		cli.Warn("Couldn't fetch versions.json. Using latest client library version.")
		return "latest"
	}
	defer res.Body.Close()

	var versions map[string]string
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		cli.Warn("Invalid versions.json. Using latest client library version.")
		return "latest"
	}

	return CompatibleVersion(versions, cgVersion)
}
