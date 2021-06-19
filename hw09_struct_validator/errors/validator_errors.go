package errors

import (
	"errors"
	"fmt"
	"strings"
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

type ValidationError struct {
	Field string
	Err   error
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{field, err}
}

func (v ValidationError) Error() string {
	var msgBuilder strings.Builder

	msgBuilder.WriteString(v.Field)

	var err ValidationErrors
	if errors.As(v.Err, &err) {
		for _, e := range err {
			msgBuilder.WriteString(fmt.Sprintf("\n\t%v", e.Error()))
		}
	} else {
		msgBuilder.WriteString(fmt.Sprintf(" %v", err.Error()))
	}
	return msgBuilder.String()
}

type ValidationErrors []*ValidationError

func (v ValidationErrors) Error() string {
	var msgBuilder strings.Builder
	for _, e := range v {
		msgBuilder.WriteString(e.Error())
	}
	return msgBuilder.String()
}
