package store

import (
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgresDB is a Postgres database.
type postgresDB gorm.DB

// NewPostgresDB creates a new Postgres database.
func newPostgresDB() *postgresDB {
	return &postgresDB{}
}

// getDSN returns the data source name.
func (pgdb *postgresDB) getDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", settings.Host, settings.Port, settings.User, settings.Password)
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.Host, settings.Port, settings.User, settings.Password, settings.DBName)
}

// newDB creates a new database.
func (pgdb *postgresDB) newDB(settings *config.DBSettings, dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// createDBIfNotExist creates a database if it does not exist.
func (pgdb *postgresDB) createDBIfNotExist(db *gorm.DB, settings *config.DBSettings) error {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", settings.DBName)
	return db.Exec(sql).Error
}
