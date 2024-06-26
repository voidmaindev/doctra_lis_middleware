package store

import (
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mysqlDB is a MySQL database.
type mysqlDB gorm.DB

// NewMySQLDB creates a new MySQL database.
func newMySQLDB() *mysqlDB {
	return &mysqlDB{}
}

// getDSN returns the data source name.
func (mysqldb *mysqlDB) getDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)", settings.User, settings.Password, settings.Host, settings.Port)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}

// newDB creates a new database.
func (mysqldb *mysqlDB) newDB(settings *config.DBSettings, dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// createDBIfNotExist creates a database if it does not exist.
func (mysqldb *mysqlDB) createDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", settings.DBName)
	return db.Exec(sql).Error
}