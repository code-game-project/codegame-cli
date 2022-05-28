package cli

import (
	"errors"
	"fmt"
)

func Info(format string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf(format, a...))
}

func Begin(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Finish() {
	fmt.Print(" done.\n")
}

func Success(format string, a ...any) {
	fmt.Printf("\x1b[32m%s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	fmt.Printf("\x1b[33mWARNING: %s\x1b[0m\n", fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) error {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\x1b[1;31mERROR: %s\x1b[0m\n", message)
	return errors.New(message)
}
