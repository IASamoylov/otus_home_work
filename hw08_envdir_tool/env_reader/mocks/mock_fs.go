package mocks

import (
	"os"

	"github.com/spf13/afero"
)

type mockFs struct {
	fs afero.Fs
}

type dirEntry struct {
	fileInfo os.FileInfo
}

func NewMockFS(fs afero.Fs) *mockFs {
	return &mockFs{fs}
}

func (mockFs *mockFs) ReadDir(name string) ([]os.DirEntry, error) {
	files, err := afero.ReadDir(mockFs.fs, name)
	entries := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		entries = append(entries, &dirEntry{fileInfo: file})
	}
	return entries, err
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
