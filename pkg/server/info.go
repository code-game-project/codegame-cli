package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/code-game-project/codegame-cli/pkg/external"
)

type GameInfo struct {
	Name          string `json:"name"`
	CGVersion     string `json:"cg_version"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	RepositoryURL string `json:"repository_url"`
}

func (a *API) FetchGameInfo() (GameInfo, error) {
	url := a.baseURL + "/info"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return GameInfo{}, fmt.Errorf("Couldn't access %s.", url)
	}
	if !external.HasContentType(res.Header, "application/json") {
		return GameInfo{}, fmt.Errorf("%s doesn't return JSON.", url)
	}
	defer res.Body.Close()

	var data GameInfo
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil || data.Name == "" || data.CGVersion == "" {
		return GameInfo{}, fmt.Errorf("Couldn't decode /info data.")
	}

	return data, nil
}
