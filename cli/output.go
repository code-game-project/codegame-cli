package cli

import (
	"errors"
	"fmt"

	"github.com/mattn/go-colorable"
)

var inProgress = false

var out = colorable.NewColorableStdout()

func IsInProgress() bool {
	return inProgress
}

func Info(format string, a ...any) {
	if inProgress {
		fmt.Fprintln(out)
		inProgress = false
	}
	fmt.Fprintf(out, "%s\n", fmt.Sprintf(format, a...))
}

func Begin(format string, a ...any) {
	fmt.Fprintf(out, format, a...)
	inProgress = true
}

func Finish() {
	if inProgress {
		fmt.Fprint(out, " done.\n")
		inProgress = false
	}
}

func Success(format string, a ...any) {
	if inProgress {
		fmt.Fprintln(out)
		inProgress = false
	}
	fmt.Fprintf(out, "\x1b[32m%s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	if inProgress {
		fmt.Fprintln(out)
		inProgress = false
	}
	fmt.Fprintf(out, "\x1b[33mWARNING: %s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) error {
	if inProgress {
		fmt.Fprintln(out)
		inProgress = false
	}
	message := fmt.Sprintf(format, a...)
	fmt.Fprintf(out, "\x1b[1;31mERROR: %s\x1b[0m\n", message)
	return errors.New(message)
}
