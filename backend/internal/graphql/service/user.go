package service

import (
	"github.com/quansolashi/golang-boierplate/backend/ent"
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
)

func NewUser(user *entity.User) *ent.User {
	return &ent.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func NewUsers(users entity.Users) []*ent.User {
	res := make([]*ent.User, len(users))
	for i := range users {
		res[i] = NewUser(users[i])
	}
	return res
}
