package errs

import "errors"

var (
	ErrinternalError = errors.New("unexpected internal error")
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotFound      = errors.New("resource not found")
)
