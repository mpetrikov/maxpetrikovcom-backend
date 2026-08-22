package labexecution

import "errors"

type TerminalError struct {
	err error
}

func NewTerminalError(err error) error {
	if err == nil {
		return nil
	}

	return &TerminalError{
		err: err,
	}
}

func (e *TerminalError) Error() string {
	return e.err.Error()
}

func (e *TerminalError) Unwrap() error {
	return e.err
}

func IsTerminalError(err error) bool {
	var terminalError *TerminalError

	return errors.As(err, &terminalError)
}
