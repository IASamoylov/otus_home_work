package common

import (
	"io"
	"os"
)

//go:generate mockgen -destination=mocks/mock_os.go -package=mocks . OSFunctions
type OSWrapper struct{}

type OSFunctions interface {
	ReadDir(name string) ([]os.DirEntry, error)
	Open(path string) (io.Reader, error)
	Setenv(key, value string) error
	Unsetenv(key string) error
	Environ() []string
}

func (OSWrapper) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSWrapper) Open(name string) (io.Reader, error) {
	return os.Open(name)
}

func (OSWrapper) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (OSWrapper) Unsetenv(key string) error {
	return os.Unsetenv(key)
}

func (OSWrapper) Environ() []string {
	return os.Environ()
}
