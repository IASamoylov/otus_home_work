package envreader

import (
	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/common"
)

//go:generate mockgen -destination mocks/mock_fs.go -package=mocks --build_flags=--mod=mod os DirEntry,FileInfo

type Ctx struct {
	os common.OSFunctions
}

func NewContext(os common.OSFunctions) *Ctx {
	return &Ctx{os}
}
