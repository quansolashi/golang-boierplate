package mysql

import (
	"context"
	"errors"
	"fmt"

	gmysql "github.com/go-sql-driver/mysql"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
	"gorm.io/gorm"
)

func NewDatabase(client *mysql.Client) *database.Database {
	return &database.Database{
		User: newUserClient(client),
	}
}

func dbError(err error) error {
	if err == nil {
		return nil
	}

	var mysqlError *gmysql.MySQLError
	if errors.As(err, &mysqlError) {
		switch mysqlError.Number {
		case 1062:
			return fmt.Errorf("%w: %s", database.ErrAlreadyExists, err.Error())
		default:
			return fmt.Errorf("%w: %s", database.ErrInternal, err.Error())
		}
	}

	switch {
	case errors.Is(err, database.ErrInvalidArgument),
		errors.Is(err, database.ErrNotFound),
		errors.Is(err, database.ErrAlreadyExists),
		errors.Is(err, database.ErrFailedPrecondition),
		errors.Is(err, database.ErrNotImplemented),
		errors.Is(err, database.ErrInternal),
		errors.Is(err, database.ErrCanceled),
		errors.Is(err, database.ErrDeadlineExceeded),
		errors.Is(err, database.ErrUnknown):
		return err // 2重ラップを防ぐため
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", database.ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", database.ErrDeadlineExceeded, err.Error())
	case errors.Is(err, gorm.ErrEmptySlice),
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrInvalidTransaction),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength),
		errors.Is(err, gorm.ErrMissingWhereClause),
		errors.Is(err, gorm.ErrModelValueRequired),
		errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return fmt.Errorf("%w: %s", database.ErrInvalidArgument, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%w: %s", database.ErrNotFound, err.Error())
	case errors.Is(err, gorm.ErrNotImplemented):
		return fmt.Errorf("%w: %s", database.ErrNotImplemented, err.Error())
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return fmt.Errorf("%w: %s", database.ErrInternal, err.Error())
	default:
		return fmt.Errorf("%w: %s", database.ErrUnknown, err.Error())
	}
}
