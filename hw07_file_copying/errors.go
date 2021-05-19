package main

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrPathTheSame           = errors.New("from and to path the same")
)

func NewErrArgumentNegative(field string) error {
	return fmt.Errorf("%s must be greater than zero", field)
}

func NewErrArgumentPath(field string) error {
	return fmt.Errorf("%s path is invalid", field)
}
