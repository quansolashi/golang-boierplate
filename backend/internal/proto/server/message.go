package server

import (
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	messagepb "github.com/quansolashi/golang-boierplate/backend/pkg/proto"
)

type message struct {
	messagepb.UnimplementedMessageServiceServer
	db *database.Database
}
