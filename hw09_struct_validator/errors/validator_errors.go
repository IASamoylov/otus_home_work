package errors

import (
	"fmt"
)

type ValidatorError struct {
	Msg string
	Err error
}

func NewValidatorError(msg string) *ValidatorError {
	return &ValidatorError{Msg: msg}
}

func NewValidatorErrorF(msg string, a ...interface{}) *ValidatorError {
	return &ValidatorError{Msg: fmt.Sprintf(msg, a...)}
}

func NewValidatorErrorW(msg string, err error) *ValidatorError {
	return &ValidatorError{msg, err}
}

func NewValidatorErrorWF(msg string, err error, a ...interface{}) *ValidatorError {
	return &ValidatorError{fmt.Sprintf(msg, a...), err}
}

func (v ValidatorError) Error() string {
	return v.Msg
}

func (v ValidatorError) Unwrap() error {
	return v.Err
}
