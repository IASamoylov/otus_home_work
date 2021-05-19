package main

import (
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"
)

type Context struct {
	source            io.Reader
	destination       io.Writer
	sourceCloser      func()
	destinationCloser func()
}

type File struct {
	stat Stat
	file *os.File
}

type Stat struct {
	path  string
	name  string
	size  int64
	isDir bool
}

type OnlyWriter struct {
	io.Writer
}

// CopyProgressBar - preset without speed and any timers. Only counters, bar and percents.
// Example: 'Prefix 20/100 [-->______] 20% Suffix'.
var CopyProgressBar pb.ProgressBarTemplate = `Copies {{string . "fileName" | green}} with offset {{string . "offset" | blue}} {{counters . }} {{bar . }} {{percent . }}`

// Copy copies the file / files to the specified directory or specified file
// can copy from any position in the file with read-limited bytes.
func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validateArguments(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	stat, err := getStat(fromPath)
	if err != nil {
		return err
	}

	if stat.isDir {
		return copyDirToFile(stat, toPath, offset, limit)
	}

	return copyFileToFile(stat, toPath, offset, limit)
}

// copyDirToFile copies all files from the directory to specified file.
// a progress bar will be visible only on big files because of
// io.Copy has default buffer on 32768 byes that greater than test files size.
func copyDirToFile(fromDirStat Stat, toPath string, offset, limit int64) error {
	entries, err := os.ReadDir(fromDirStat.path)
	if err != nil {
		return err
	}

	destination, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	defer func() {
		destination.Close()
		destination.Sync()
	}()

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fromFilePath := path.Join(fromDirStat.path, entry.Name())

		stat, err := getStat(fromFilePath)
		if err != nil {
			return err
		}

		if stat.size == 0 {
			return ErrUnsupportedFile
		}

		source, err := os.OpenFile(fromFilePath, os.O_RDONLY, 0)
		if err != nil {
			return err
		}

		context, err := prepareContext(File{stat: stat, file: source}, File{file: destination}, offset, limit)
		if err != nil {
			return err
		}

		defer context.sourceCloser()

		_, err = io.Copy(context.destination, context.source)

		return err
	}

	return nil
}

// copyFileToFile copies the file from the directory to specified directory or specified file.
func copyFileToFile(stat Stat, toPath string, offset, limit int64) error {
	source, err := os.OpenFile(stat.path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	if stat.size == 0 {
		return ErrUnsupportedFile
	}

	destination, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}

	context, err := prepareContext(File{stat: stat, file: source}, File{file: destination}, offset, limit)
	if err != nil {
		return err
	}

	defer context.destinationCloser()
	defer context.sourceCloser()

	_, err = io.Copy(context.destination, context.source)

	return err
}

// prepareContext returns context of copying.
func prepareContext(source, destination File, offset, limit int64) (*Context, error) {
	if source.stat.size < offset {
		return nil, ErrOffsetExceedsFileSize
	}

	fileSize := source.stat.size

	if limit == 0 {
		limit = fileSize
	}

	if limit > fileSize-offset {
		limit = fileSize - offset
	}
	_, err := source.file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	bar := CopyProgressBar.Start64(limit)
	bar.Set("fileName", source.stat.name)
	bar.Set("offset", bar.Format(offset))

	reader := io.LimitReader(source.file, limit)
	barReader := bar.NewProxyReader(reader)

	return &Context{
		source:      barReader,
		destination: OnlyWriter{destination.file},
		sourceCloser: func() {
			barReader.Close()
		},
		destinationCloser: func() {
			destination.file.Close()
			destination.file.Sync()
		},
	}, nil
}

// getStat get file stat by path.
func getStat(path string) (Stat, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return Stat{}, err
	}

	return Stat{
		path:  path,
		name:  stat.Name(),
		isDir: stat.IsDir(),
		size:  stat.Size(),
	}, nil
}

// validateArguments validates input parameters that require for coping.
func validateArguments(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return NewErrArgumentPath("from")
	}

	if toPath == "" {
		return NewErrArgumentPath("to")
	}

	if toPath == fromPath {
		return ErrPathTheSame
	}

	if offset < 0 {
		return NewErrArgumentNegative("offset")
	}

	if limit < 0 {
		return NewErrArgumentNegative("limit")
	}

	return nil
}
