package cmd

import (
	"errors"
	"fmt"
	"os"

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

func init() {}

// abort prints the error to the console and terminates the program.
// abort does nothing if err is nil.
func abort(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, cli.ErrCanceled) {
		cli.Print(err.Error())
	} else {
		cli.Error(err.Error())
	}

	if err == cli.ErrCanceled {
		os.Exit(0)
	}
	os.Exit(1)
}

// abortf prints the error to the console and terminates the program.
// abortf does nothing if err is nil.
func abortf(format string, err error) {
	if err == nil {
		return
	}

	if errors.Is(err, cli.ErrCanceled) {
		cli.Print(format, err)
	} else {
		cli.Error(fmt.Errorf(format, err).Error())
	}

	if err == cli.ErrCanceled {
		os.Exit(0)
	}
	os.Exit(1)
}
