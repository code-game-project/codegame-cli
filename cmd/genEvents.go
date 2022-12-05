package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/code-game-project/go-utils/cgfile"
	"github.com/code-game-project/go-utils/cggenevents"
	"github.com/code-game-project/go-utils/exec"
	"github.com/spf13/cobra"
)

// genEventsCmd represents the genEvents command
var genEventsCmd = &cobra.Command{
	Use:   "gen-events",
	Short: "Generate event definitions from CGE files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var output string
		output, err := cmd.Flags().GetString("output")
		abort(err)
		languages, err := cmd.Flags().GetStringSlice("languages")
		abort(err)

		var filename string
		if len(args) == 0 {
			root, err := cgfile.FindProjectRootRelative()
			if err != nil || cmd.Flags().Changed("output") || cmd.Flags().Changed("languages") {
				fmt.Println(err, output, languages)
				abort(errors.New("Expected game URL."))
			}
			data, err := cgfile.LoadCodeGameFile(root)
			if err != nil {
				abortf("Failed to load CodeGame data: %s", err)
			}
			if data.Type == "client" {
				abort(errors.New("Use `codegame update` instead."))
			} else if data.Type != "server" {
				fmt.Println("hi")
				abort(errors.New("Expected game URL."))
			} else {
				filename = filepath.Join(root, "events.cge")
				languages = []string{data.Lang}
				switch data.Lang {
				case "go":
					output = filepath.Join(root, strings.ReplaceAll(strings.ReplaceAll(data.Game, "_", ""), "-", ""))
				default:
					abort(errors.New("Expected game URL."))
					fmt.Println("hi2")
				}
			}
		}

		var cge []byte
		if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
			if !strings.HasSuffix(filename, "/api/events") && !strings.HasSuffix(filename, ".cge") {
				if strings.HasSuffix(filename, "/api") {
					filename += "/events"
				} else if strings.HasSuffix(filename, "/") {
					filename += "api/events"
				} else {
					filename += "/api/events"
				}
			}
			resp, err := http.Get(filename)
			if err != nil {
				abort(fmt.Errorf("Failed to reach url '%s': %s", filename, err))
			}
			if resp.StatusCode != http.StatusOK {
				abort(fmt.Errorf("Failed to download CGE file from url '%s'", filename))
			}
			if !strings.Contains(resp.Header.Get("Content-Type"), "text/plain") {
				abort(fmt.Errorf("Unsupported content type at '%s': expected %s, got %s\n", filename, "text/plain", resp.Header.Get("Content-Type")))
			}
			cge, err = io.ReadAll(resp.Body)
			abortf("Failed to read CGE file: %s", err)
		} else {
			cge, err = os.ReadFile(filename)
			abortf("Failed to read CGE file: %s", err)
		}

		cgeVersion, err := cggenevents.ParseCGEVersion(string(cge))
		abortf("Failed to determine CGE file version: %s", err)

		cgGenEvents, err := cggenevents.InstallCGGenEvents(cgeVersion)
		abortf("Failed to install cg-gen-events: %s", err)
		_, err = exec.Execute(false, filepath.Join(xdg.DataHome, "codegame", "bin", "cg-gen-events", cgGenEvents), filename, "-o", output, "-l", strings.Join(languages, ","))
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(genEventsCmd)
	genEventsCmd.Flags().StringP("output", "o", ".", "The directory where every file will be generated into. (Will be created if it does not exist.)")
	genEventsCmd.Flags().StringSliceP("languages", "l", []string{""}, "A list of target languages.")
}
