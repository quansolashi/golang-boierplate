package server

import (
	"context"
	"time"

	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	userpb "github.com/quansolashi/golang-boierplate/backend/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Params struct {
	DB *database.Database
}

type Server struct {
	Grpc     *grpc.Server
	Services *GrpcServices
}

type GrpcServices struct {
	User User
}

type User interface {
	GetUser(ctx context.Context, r *userpb.UserRequest) (*userpb.User, error)
}

func NewGrpcServer(params *Params) *Server {
	grpc := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Timeout:           10 * time.Second,
		}),
	)

	// debugging
	reflection.Register(grpc)

	return &Server{
		Grpc: grpc,
		Services: &GrpcServices{
			User: newUserService(params.DB, grpc),
		},
	}
}
