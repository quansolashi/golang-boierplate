package mysql

import (
	"github.com/quansolashi/message-extractor/backend/internal/database"
	"github.com/quansolashi/message-extractor/backend/pkg/mysql"
)

func NewDatabase(client *mysql.Client) *database.Database {
	return &database.Database{}
}
