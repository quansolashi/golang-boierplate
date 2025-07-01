package service

import (
	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/internal/web/response"
)

type User struct {
	response.User
}

type Users []*User

func NewUser(user *entity.User) *User {
	return &User{
		User: response.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
}

func (u *User) Response() *response.User {
	return &u.User
}

func NewUsers(users []*entity.User) Users {
	res := make([]*User, len(users))
	for i := range users {
		res[i] = NewUser(users[i])
	}
	return res
}

func (us Users) Response() response.Users {
	res := make([]*response.User, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}
