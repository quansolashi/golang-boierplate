package server

import (
	"context"

	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	userpb "github.com/quansolashi/golang-boierplate/backend/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type user struct {
	userpb.UnimplementedUserServiceServer
	db *database.Database
}

func newUserService(db *database.Database, grpc *grpc.Server) User {
	s := &user{
		db: db,
	}
	userpb.RegisterUserServiceServer(grpc, s)
	return s
}

func (u *user) GetUser(ctx context.Context, r *userpb.UserRequest) (*userpb.User, error) {
	userID := r.GetId()
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field: ID")
	}

	user, err := u.db.User.Get(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userpb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
