package util

import "errors"

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrNotAcceptable      = errors.New("not acceptable")
	ErrAlreadyExists      = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrResourceExhausted  = errors.New("resource exhausted")
	ErrNotImplemented     = errors.New("not implemented")
	ErrInternal           = errors.New("internal error")
	ErrCanceled           = errors.New("canceled")
	ErrUnavailable        = errors.New("unavailable")
	ErrDeadlineExceeded   = errors.New("deadline exceeded")
	ErrUnknown            = errors.New("unknown")
)
