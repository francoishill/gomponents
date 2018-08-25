package clienterror

import (
	"github.com/pkg/errors"
)

//Error holds a Status int for usage in HTTP
type Error interface {
	error
	Status() int
}

func NewError(err error, status int) Error {
	return &localError{
		Err:    errors.WithStack(err),
		status: status,
	}
}

func NewErrorDefaultStatus(err error) Error {
	return NewError(err, 500)
}

type localError struct {
	Err    error
	status int
}

func (e *localError) Error() string { return e.Err.Error() }
func (e *localError) Status() int   { return e.status }
