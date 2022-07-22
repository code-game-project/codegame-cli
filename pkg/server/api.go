package server

import (
	"fmt"
	"net/http"

	"github.com/code-game-project/codegame-cli/pkg/external"
)

type API struct {
	baseURL string
}

func NewAPI(url string) (*API, error) {
	url = external.TrimURL(url)
	tls := external.IsTLS(url)

	api := &API{
		baseURL: external.BaseURL("http", tls, url),
	}

	resp, err := http.Get(api.baseURL + "/api/info")
	if err != nil {
		return nil, fmt.Errorf("Cannot reach %s.", api.baseURL)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK && external.HasContentType(resp.Header, "application/json") {
		api.baseURL += "/api"
	}

	return api, nil
}

func (a *API) BaseURL() string {
	return a.baseURL
}
