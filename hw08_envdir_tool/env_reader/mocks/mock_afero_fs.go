package mocks

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/afero"
)

type MockFs struct {
	fs afero.Fs
}

type dirEntry struct {
	fileInfo os.FileInfo
}

func NewAferoWrapper(fs afero.Fs) *MockFs {
	return &MockFs{fs}
}

func (mockFs *MockFs) ReadDir(name string) ([]os.DirEntry, error) {
	files, err := afero.ReadDir(mockFs.fs, name)
	entries := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		entries = append(entries, &dirEntry{fileInfo: file})
	}
	return entries, err
}

func (mockFs *MockFs) Open(path string) (io.Reader, error) {
	return mockFs.fs.Open(path)
}

func (mockFs *MockFs) Setenv(key, value string) error {
	return errors.New("not implemented")
}

func (mockFs *MockFs) Unsetenv(key string) error {
	return errors.New("not implemented")
}

func (mockFs *MockFs) Environ() []string {
	return nil
}

func (e *dirEntry) Name() string {
	return e.fileInfo.Name()
}

func (e *dirEntry) IsDir() bool {
	return e.fileInfo.IsDir()
}

func (e *dirEntry) Type() os.FileMode {
	return e.fileInfo.Mode()
}

func (e *dirEntry) Info() (os.FileInfo, error) {
	return e.fileInfo, nil
}
