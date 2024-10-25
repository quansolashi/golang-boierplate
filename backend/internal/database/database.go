package database

import (
	"context"
	"errors"

	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
)

var (
	ErrInvalidArgument    = errors.New("database: invalid argument")
	ErrNotFound           = errors.New("database: not found")
	ErrNotAcceptable      = errors.New("database: not acceptable")
	ErrAlreadyExists      = errors.New("database: already exists")
	ErrFailedPrecondition = errors.New("database: failed precondition")
	ErrNotImplemented     = errors.New("database: not implemented")
	ErrInternal           = errors.New("database: internal error")
	ErrCanceled           = errors.New("database: canceled")
	ErrDeadlineExceeded   = errors.New("database: deadline exceeded")
	ErrUnknown            = errors.New("database: unknown")
)

type Database struct {
	User User
}

type User interface {
	List(ctx context.Context) (entity.Users, error)
	Get(ctx context.Context, userID uint64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
