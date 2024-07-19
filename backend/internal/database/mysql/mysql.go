package mysql

import (
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
)

func NewDatabase(client *mysql.Client) *database.Database {
	return &database.Database{
		User: newUserClient(client),
	}
}
