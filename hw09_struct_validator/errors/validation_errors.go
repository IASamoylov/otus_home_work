package errors

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{field, err}
}

func (v ValidationError) Error() string {
	var msgBuilder strings.Builder

	var errs *ValidationErrors
	if errors.As(v.Err, &errs) {
		for _, e := range errs.Errors {
			msgBuilder.WriteString(fmt.Sprintf("%s.%v", v.Field, e.Error()))
		}
	} else {
		msgBuilder.WriteString(fmt.Sprintf("%s: %v", v.Field, v.Err.Error()))
	}
	return msgBuilder.String()
}

type ValidationErrors struct {
	Errors []*ValidationError
}

func (v ValidationErrors) Error() string {
	var msgBuilder strings.Builder
	for _, e := range v.Errors {
		msgBuilder.WriteString(e.Error())
	}
	return msgBuilder.String()
}
