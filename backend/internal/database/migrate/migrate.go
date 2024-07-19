package migrate

import (
	"context"

	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
	"go.uber.org/zap"
)

var tables = map[string]interface{}{
	"users": &entity.User{},
}

var migrationOrders = []string{
	"users",
}

func Run(ctx context.Context, db *mysql.Client, logger *zap.Logger) error {
	// logger.Info("Start db migration...")
	for _, tableName := range migrationOrders {
		model, exists := tables[tableName]
		if !exists {
			// logger.Error("Model not found in Tables map", zap.String("table", tableName))
			continue
		}
		if err := db.DB.WithContext(ctx).Table(tableName).AutoMigrate(model); err != nil {
			// logger.Error("Failed to migrate", zap.String("table", tableName), zap.Error(err))
			return err
		}
	}
	// logger.Info("Finished db migration")
	return nil
}
