package executor

import (
	"io"

	common2 "github.com/IASamoylov/otus_home_work/hw08_envdir_tool/common"
)

type Ctx struct {
	stdIn  io.Reader
	stdErr io.Writer
	stdOut io.Writer
	os     common2.OSFunctions
}

func NewContext(os common2.OSFunctions, stdIn io.Reader, stdOut, stdErr io.Writer) *Ctx {
	return &Ctx{stdIn, stdErr, stdOut, os}
}
