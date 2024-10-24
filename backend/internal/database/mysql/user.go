package mysql

import (
	"context"

	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
)

type user struct {
	db *mysql.Client
}

func newUserClient(client *mysql.Client) database.User {
	return &user{
		db: client,
	}
}

func (u *user) List(ctx context.Context) (entity.Users, error) {
	var users entity.Users

	stmt := u.db.DB.WithContext(ctx).Table("users")

	err := stmt.Find(&users).Error
	return users, err
}

func (u *user) Get(ctx context.Context, userID uint64) (*entity.User, error) {
	var user *entity.User

	stmt := u.db.DB.WithContext(ctx).
		Table("users").
		Where("id", userID)

	err := stmt.First(&user).Error
	return user, err
}
