package envreader

import (
	"io"
	"os"
)

//go:generate mockgen -destination=mocks/mock_os.go -package=mocks . OS
//go:generate mockgen -destination mocks/mock_fs.go -package=mocks --build_flags=--mod=mod os DirEntry,FileInfo

type Ctx struct {
	os OS
}

type OS interface {
	ReadDir(name string) ([]os.DirEntry, error)
	Open(path string) (io.Reader, error)
}

func NewContext(os OS) *Ctx {
	return &Ctx{os}
}

func NewOSContext() *Ctx {
	return &Ctx{osWrap{}}
}

type osWrap struct{}

func (osWrap) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (osWrap) Open(name string) (io.Reader, error) {
	return os.Open(name)
}
