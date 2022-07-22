package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/code-game-project/codegame-cli/pkg/external"
)

func (a *API) GetCGEFile() (string, error) {
	res, err := http.Get(a.baseURL + "/events")
	if err != nil || res.StatusCode != http.StatusOK || (!external.HasContentType(res.Header, "text/plain") && !external.HasContentType(res.Header, "application/octet-stream")) {
		return "", fmt.Errorf("Couldn't access /events endpoint.")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Couldn't read /events file.")
	}
	return string(data), nil
}
