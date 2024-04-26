package store

import (
	"errors"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"gorm.io/gorm"
)

// DB is an interface that defines the methods that a database should implement.
type DB interface {
	getDSN(settings *config.DBSettings, woDBName bool) string
	newDB(settings *config.DBSettings, dsn string) (*gorm.DB, error)
	createDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error
}

func NewDB(settings *config.DBSettings, config *gorm.Config) (*gorm.DB, error) {
	sqlDB := getSqlDB(settings)

	if settings.CreateDB {
		dsnWoDBName := sqlDB.getDSN(settings, true)
		if dsnWoDBName == "" {
			return nil, errors.New("failed to get DSN without DB name")
		}

		dbWoDBName, err := sqlDB.newDB(settings, dsnWoDBName)
		if err != nil {
			return nil, errors.New("failed to connect to a new DB without DB name")
		}

		err = sqlDB.createDBIfNotExist(dbWoDBName, settings)
		if err != nil {
			return nil, errors.New("failed to create DB if not exist")
		}
	}

	dsn := sqlDB.getDSN(settings, false)
	if dsn == "" {
		return nil, errors.New("failed to get DSN")
	}

	return sqlDB.newDB(settings, dsn)
}

func getSqlDB(settings *config.DBSettings) DB {
	var sqlDB DB

	switch settings.DriverName {
	case "sqlserver":
		sqlDB = newSQLServerDB()
	case "postgres":
		sqlDB = newPostgresDB()
	case "mysql":
		sqlDB = newMySQLDB()
	}

	return sqlDB
}

