package database

import (
	"context"

	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
)

type Database struct {
	User User
}

type User interface {
	List(ctx context.Context) (entity.Users, error)
	Get(ctx context.Context, userID uint64) (*entity.User, error)
}
