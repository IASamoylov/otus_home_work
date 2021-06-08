package env_reader

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Name       string
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func (ctx *EnvReaderCtx) ReadDir(dir string) (Environment, error) {
	entries, err := ctx.os.ReadDir(dir)
	if err != nil {
		return nil, NewEnvReaderErr("Errors occurred while getting environment from the directory", err)
	}

	environment := make(Environment, len(entries))
	for _, entry := range entries {
		if validateEntry(entry) {
			continue
		}
		var env EnvValue
		if env, err = ctx.parseEntry(dir, entry); err != nil {
			return nil, err
		}
		environment[env.Name] = env
	}

	return environment, nil
}

// validateEntry ignores files that are a directory or have '=' in name
func validateEntry(entry os.DirEntry) bool {
	if entry.IsDir() || strings.ContainsRune(entry.Name(), '=') {
		return true
	}

	return false
}

func (ctx *EnvReaderCtx) parseEntry(dir string, entry os.DirEntry) (env EnvValue, err error) {
	info, err := entry.Info()
	if err != nil {
		return EnvValue{}, NewEnvReaderErrF("Error processing file %v", err, entry.Name())
	}

	if info.Size() == 0 {
		return EnvValue{Name: entry.Name(), NeedRemove: true}, nil
	}

	file, err := ctx.os.Open(path.Join(dir, info.Name()))
	if err != nil {
		return EnvValue{}, NewEnvReaderErrF("Error processing file %v", err, entry.Name())
	}

	reader := bufio.NewReader(file)
	data, _, err := reader.ReadLine()
	if err != nil {
		return EnvValue{}, NewEnvReaderErrF("Error processing file %v", err, entry.Name())
	}

	data = bytes.TrimFunc(data, func(r rune) bool {
		return unicode.IsSpace(r) || rune("\x00"[0]) == r
	})
	return EnvValue{
		Name:  entry.Name(),
		Value: string(data),
	}, nil
}