package executor

import "io"

type ExecutorCtx struct {
	stdIn  io.Reader
	stdErr io.Writer
	stdOut io.Writer
}

func NewExecutorCtx(stdIn io.Reader, stdOut, stdErr io.Writer) *ExecutorCtx {
	return &ExecutorCtx{stdIn, stdErr, stdOut}
}
