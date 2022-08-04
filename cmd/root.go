package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/Bananenpro/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codegame",
	Short: "The official CodeGame CLI.",
	Long:  `The CodeGame CLI helps you develop CodeGame applications.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

// abort prints the error to the console and terminates the program.
// abort does nothing if err is nil.
func abort(err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*exec.ExitError); !ok && !errors.Is(err, cli.ErrCanceled) {
		cli.Error(err.Error())
	}
	os.Exit(1)
}

// abortf prints the error to the console and terminates the program.
// abortf does nothing if err is nil.
func abortf(format string, err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*exec.ExitError); !ok && !errors.Is(err, cli.ErrCanceled) {
		cli.Error(format, err)
	}
	os.Exit(1)
}
