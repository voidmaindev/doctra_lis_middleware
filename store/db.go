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
	if settings.CreateDB {
		dsnWoDBName := getDSN(settings, true)
		if dsnWoDBName == "" {
			return nil, errors.New("failed to get DSN without DB name")
		}

		dbWoDBName, err := newDB(settings, dsnWoDBName)
		if err != nil {
			return nil, errors.New("failed to connect to a new DB without DB name")
		}

		err = createDBIfNotExist(dbWoDBName, settings)
		if err != nil {
			return nil, errors.New("failed to create DB if not exist")
		}
	}

	dsn := getDSN(settings, false)
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

func newDB(settings *config.DBSettings, dsn string) (*gorm.DB, error) {
	switch settings.DriverName {
	case "sqlserver":
		return getSQLServerDB(dsn, &gorm.Config{})
	case "postgres":
		return getPostgresDB(dsn, &gorm.Config{})
	case "mysql":
		return getMySQLDB(dsn, &gorm.Config{})
	}

	return nil, errors.New("unsupported sql driver")
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

func createDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	switch settings.DriverName {
	case "sqlserver":
		return createSQLServerDBIfNotExist(db, settings)
	case "postgres":
		return createPostgresDBIfNotExist(db, settings)
	case "mysql":
		return createMySQLDBIfNotExist(db, settings)
	}

	return errors.New("unsupported sql driver")
}

func createSQLServerDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("IF NOT EXISTS (SELECT name FROM master.dbo.sysdatabases WHERE name = N'%s') CREATE DATABASE %s", settings.DBName, settings.DBName)
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

func createPostgresDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", settings.DBName)
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

func createMySQLDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", settings.DBName)
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}
