package executor

import "fmt"

type executorErr struct {
	msg string
	err error
}

func NewExecutorErr(msg string, err error) error {
	return &executorErr{msg, err}
}

func NewExecutorErrF(msg string, err error, a ...interface{}) error {
	return &executorErr{fmt.Sprintf(msg, a...), err}
}

func (err *executorErr) Error() string {
	return err.msg
}

func (err *executorErr) Unwrap() error {
	return err.err
}
