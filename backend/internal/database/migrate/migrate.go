package migrate

import (
	"context"
	"fmt"

	"github.com/quansolashi/golang-boierplate/backend/internal/entity"
	"github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
	"go.uber.org/zap"
)

var tables = map[string]interface{}{
	"users":         &entity.User{},
	"user_profiles": &entity.UserProfile{},
	"posts":         &entity.Post{},
}

var migrationOrders = []string{
	"users",
	"user_profiles",
	"posts",
}

func Run(ctx context.Context, db *mysql.Client, logger *zap.Logger) error {
	// temporary use logger to avoid rules of lint
	fmt.Println(logger.Name())

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
