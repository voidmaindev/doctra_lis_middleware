package store

import (
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// sqlserverDB is a SQL Server database.
type sqlserverDB gorm.DB

// NewSQLServerDB creates a new SQL Server database.
func newSQLServerDB() *sqlserverDB {
	return &sqlserverDB{}
}

// getDSN returns the data source name.
func (sqldb *sqlserverDB) getDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s", settings.User, settings.Password, settings.Host, settings.Port)
	}

	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}

// newDB creates a new database.
func (sqldb *sqlserverDB) newDB(settings *config.DBSettings, dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}

// createDBIfNotExist creates a database if it does not exist.
func (sqldb *sqlserverDB) createDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("IF NOT EXISTS (SELECT name FROM master.dbo.sysdatabases WHERE name = N'%s') CREATE DATABASE %s", settings.DBName, settings.DBName)
	return db.Exec(sql).Error
}
