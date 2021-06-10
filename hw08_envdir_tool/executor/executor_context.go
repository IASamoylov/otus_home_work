package executor

import "io"

type Ctx struct {
	stdIn  io.Reader
	stdErr io.Writer
	stdOut io.Writer
}

func NewExecutorCtx(stdIn io.Reader, stdOut, stdErr io.Writer) *Ctx {
	return &Ctx{stdIn, stdErr, stdOut}
}
