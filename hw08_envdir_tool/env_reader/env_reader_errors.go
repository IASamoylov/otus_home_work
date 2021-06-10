package envreader

import "fmt"

type envReaderErr struct {
	msg string
	err error
}

func NewEnvReaderErr(msg string, err error) error {
	return &envReaderErr{msg, err}
}

func NewEnvReaderErrF(msg string, err error, a ...interface{}) error {
	return &envReaderErr{fmt.Sprintf(msg, a), err}
}

func (err *envReaderErr) Error() string {
	return err.msg
}

func (err *envReaderErr) Unwrap() error {
	return err.err
}
