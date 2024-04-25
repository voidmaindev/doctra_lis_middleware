package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func NewDB(settings *config.DBSettings, config *gorm.Config) (*gorm.DB, error) {
	dsn := getDSN(settings)
	if dsn == "" {
		return nil, errors.New("failed to get DSN")
	}

	switch settings.DriverName {
	case "sqlserver":
		return getSQLServerDB(dsn, config)
	case "postgres":
		return getPostgresDB(dsn, config)
	case "mysql":
		return getMySQLDB(dsn, config)
	}

	return nil, errors.New("unsupported sql driver")
}

func getDSN(settings *config.DBSettings) string {
	switch settings.DriverName {
	case "sqlserver":
		return getSQLServerDSN(settings)
	case "postgres":
		return getPostgresDSN(settings)
	case "mysql":
		return getMySQLDSN(settings)
	}

	return ""
}

func getSQLServerDSN(settings *config.DBSettings) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}

func getPostgresDSN(settings *config.DBSettings) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.Host, settings.Port, settings.User, settings.Password, settings.DBName)
}

func getMySQLDSN(settings *config.DBSettings) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}

func getSQLServerDB(dsn string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getPostgresDB(dsn string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getMySQLDB(dsn string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}
