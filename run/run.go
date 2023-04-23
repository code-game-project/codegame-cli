package run

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/code-game-project/cli-utils/cgfile"
	"github.com/code-game-project/cli-utils/cli"
	"github.com/code-game-project/cli-utils/modules"
	"github.com/code-game-project/cli-utils/server"
	"github.com/code-game-project/cli-utils/sessions"
)

func RunServer(cgFile *cgfile.CodeGameFileData, port int, args []string) error {
	mod, err := modules.LoadModule(cgFile.Language)
	if err != nil {
		return fmt.Errorf("load %s module: %w", cgFile.Language, err)
	}

	var p *int32
	if port != 0 {
		temp := int32(port)
		p = &temp
	}

	return mod.ExecRunServer(cgFile.ModVersion, cgFile.Language, p, args)
}

func RunClient(cgFile *cgfile.CodeGameFileData, spectate bool, args []string) error {
	sessns, err := sessions.ListSessionsByGame(cgFile.GameURL)
	if err != nil {
		return fmt.Errorf("list sessions: %w", err)
	}

	games := make([]string, 0, len(sessns)+2)
	games = append(games, "Create new...")
	for _, s := range sessns {
		games = append(games, s.GameID)
	}
	games = append(games, "Other...")
	index := cli.Select("Select game:", games)
	var gameID string
	var joinSecret string
	var created bool
	if index == 0 {
		public := cli.YesNo("Public?", false)
		protected := cli.YesNo("Protected?", false)
		var configData json.RawMessage
		cli.Input("Config (json):", false, "null", func(input interface{}) error {
			err = json.Unmarshal([]byte(input.(string)), &configData)
			if err != nil {
				return fmt.Errorf("Invalid JSON: %w", err)
			}
			return nil
		})
		gameID, joinSecret, err = server.CreateGame(cgFile.GameURL, public, protected, configData)
		if err != nil {
			return err
		}
		created = true
	} else if index == len(games)-1 {
		gameID = cli.Input("Game ID:", true, "")
	} else {
		gameID = sessns[index].GameID
	}

	if spectate {
		return RunClientSpectate(cgFile, gameID, args)
	}

	var playerID string
	var playerSecret string
	if !created {
		type player struct {
			id       string
			secret   string
			username string
		}
		players := make([]player, 0, 1)
		for _, s := range sessns {
			if s.GameID == gameID {
				players = append(players, player{
					id:       s.PlayerID,
					secret:   s.PlayerSecret,
					username: s.Username,
				})
			}
		}
		options := make([]string, 0, len(players)+2)
		options = append(options, "Create new...")
		for _, p := range players {
			options = append(options, fmt.Sprintf("%s (%s)", p.username, p.id))
		}
		options = append(options, "Other...")
		index := cli.Select("Select player:", options)
		if index == 0 {
			if joinSecret == "" {
				var g server.Game
				g, err = server.FetchGame(cgFile.GameURL, gameID)
				if err != nil {
					return fmt.Errorf("fetch game: %w", err)
				}
				if g.Protected {
					joinSecret = cli.Input("Join secret:", true, "")
				}
			}
			playerID, playerSecret, err = server.CreatePlayer(cgFile.GameURL, gameID, cli.Input("Username:", true, ""), joinSecret)
			if err != nil {
				return err
			}
		} else if index == len(options)-1 {
			playerID = cli.Input("Player ID:", true, "")
			playerSecret = cli.Input("Player secret:", true, "")
		} else {
			playerID = players[index].id
			playerSecret = players[index].secret
		}
	} else {
		playerID, playerSecret, err = server.CreatePlayer(cgFile.GameURL, gameID, cli.Input("Username:", true, ""), joinSecret)
		if err != nil {
			return err
		}
	}
	return RunClientConnect(cgFile, gameID, playerID, playerSecret, args)
}

func RunClientCreate(cgFile *cgfile.CodeGameFileData, username string, public, protected, spectate bool, configPath string, args []string) error {
	var configData json.RawMessage
	if configPath != "" {
		config, err := os.Open(configPath)
		if err != nil {
			return fmt.Errorf("open config file: %w", err)
		}
		defer config.Close()
		err = json.NewDecoder(config).Decode(&configData)
		if err != nil {
			return fmt.Errorf("decode config file: %w", err)
		}
	}

	gameID, joinSecret, err := server.CreateGame(cgFile.GameURL, public, protected, configData)
	if err != nil {
		return err
	}
	cli.Print("Game ID: %s", gameID)
	if joinSecret != "" {
		cli.Print("Join secret: %s", joinSecret)
	}
	if spectate {
		return RunClientSpectate(cgFile, gameID, args)
	}
	return RunClientJoin(cgFile, gameID, username, joinSecret, args)
}

func RunClientJoin(cgFile *cgfile.CodeGameFileData, gameID, username, joinSecret string, args []string) error {
	playerID, playerSecret, err := server.CreatePlayer(cgFile.GameURL, gameID, username, joinSecret)
	if err != nil {
		return err
	}
	return RunClientConnect(cgFile, gameID, playerID, playerSecret, args)
}

func RunClientConnect(cgFile *cgfile.CodeGameFileData, gameID, playerID, playerSecret string, args []string) error {
	mod, err := modules.LoadModule(cgFile.Language)
	if err != nil {
		return fmt.Errorf("load %s module: %w", cgFile.Language, err)
	}
	return mod.ExecRunClient(cgFile.ModVersion, cgFile.GameURL, cgFile.Language, gameID, &playerID, &playerSecret, false, args)
}

func RunClientSpectate(cgFile *cgfile.CodeGameFileData, gameID string, args []string) error {
	mod, err := modules.LoadModule(cgFile.Language)
	if err != nil {
		return fmt.Errorf("load %s module: %w", cgFile.Language, err)
	}
	return mod.ExecRunClient(cgFile.ModVersion, cgFile.GameURL, cgFile.Language, gameID, nil, nil, true, args)
}
