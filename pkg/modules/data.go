package modules

import (
	"encoding/json"
	"errors"
	"os"
)

type NewClientData struct {
	// Lang contains the chosen programming language (in case one module supports multiple languages).
	Lang string `json:"lang"`
	// Name contains the name of the game.
	Name string `json:"name"`
	// URL contains the URL of the game server.
	URL string `json:"url"`
	// LibraryVersion contains the version of the client library to use. (e.g. 0.9.2)
	LibraryVersion string `json:"library_version"`
}

type NewServerData struct {
	// Lang contains the chosen programming language (in case one module supports multiple languages).
	Lang string `json:"lang"`
	// LibraryVersion contains the version of the server library to use. (e.g. 0.9.2)
	LibraryVersion string `json:"library_version"`
}

type UpdateData struct {
	// Lang contains the chosen programming language (in case one module supports multiple languages).
	Lang string `json:"lang"`
	// LibraryVersion contains the new version of the library to use. (e.g. 0.9.2)
	LibraryVersion string `json:"library_version"`
}

type RunData struct {
	// Lang contains the chosen programming language (in case one module supports multiple languages).
	Lang string `json:"lang"`
	// Args contains a list of command line arguments for the application.
	Args []string `json:"args"`
}

type BuildData struct {
	// Lang contains the chosen programming language (in case one module supports multiple languages).
	Lang string `json:"lang"`
	// Output contains the output file name.
	Output string `json:"output"`
}

func ReadCommandConfig[T any]() (T, error) {
	var data T

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		return data, errors.New("empty CONFIG_FILE environment variable")
	}

	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	return data, err
}
