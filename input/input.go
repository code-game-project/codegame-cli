package input

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	ErrCanceled = errors.New("Canceled.")
)

func Input(question string) (string, error) {
	var result string
	err := survey.AskOne(&survey.Input{
		Message: question,
	}, &result, survey.WithValidator(survey.Required))
	if err == terminal.InterruptErr {
		err = ErrCanceled
	}
	return result, err
}

func YesNo(question string, defaultValue bool) (yes bool, err error) {
	err = survey.AskOne(&survey.Confirm{
		Message: question,
		Default: defaultValue,
	}, &yes, survey.WithValidator(survey.Required))
	if err == terminal.InterruptErr {
		err = ErrCanceled
	}
	return yes, err
}

func Select(msg string, displayOptions []string, options []string) (string, error) {
	var index int
	err := survey.AskOne(&survey.Select{
		Message: msg,
		Options: displayOptions,
	}, &index, survey.WithValidator(survey.Required))
	if err == terminal.InterruptErr {
		err = ErrCanceled
	}
	return options[index], err
}
