package cli

import (
	"errors"
	"fmt"
)

var inProgress = false

func IsInProgress() bool {
	return inProgress
}

func Info(format string, a ...any) {
	if inProgress {
		fmt.Println()
		inProgress = false
	}
	fmt.Printf("%s\n", fmt.Sprintf(format, a...))
}

func Begin(format string, a ...any) {
	fmt.Printf(format, a...)
	inProgress = true
}

func Finish() {
	if inProgress {
		fmt.Print(" done.\n")
		inProgress = false
	}
}

func Success(format string, a ...any) {
	if inProgress {
		fmt.Println()
		inProgress = false
	}
	fmt.Printf("\x1b[32m%s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	if inProgress {
		fmt.Println()
		inProgress = false
	}
	fmt.Printf("\x1b[33mWARNING: %s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) error {
	if inProgress {
		fmt.Println()
		inProgress = false
	}
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\x1b[1;31mERROR: %s\x1b[0m\n", message)
	return errors.New(message)
}
