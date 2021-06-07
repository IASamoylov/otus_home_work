package env_reader

import (
	"io"
	"os"
)

//go:generate mockgen -destination=mocks/mock_os_wrapper.go -package=mocks . FS

type EnvReaderCtx struct {
	os OS
}

type OSWrap struct {
}

type OS interface {
	ReadDir(name string) ([]os.DirEntry, error)
	//ReadFile(name string) ([]byte, error)
	GetReader(path string) (io.Reader, error)
}

func NewContext(os OS) *EnvReaderCtx {
	return &EnvReaderCtx{os}
}

func NewOSContext() *EnvReaderCtx {
	return &EnvReaderCtx{OSWrap{}}
}

func (osw OSWrap) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (osw OSWrap) GetReader(name string) ([]byte, error) {
	return os.Open(name)
}
