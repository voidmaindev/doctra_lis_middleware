package store

import (
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/config"
)

func getDSN(settings *config.DBSettings, woDBName bool) string {
	switch settings.DriverName {
	case "sqlserver":
		return getSQLServerDSN(settings, woDBName)
	case "postgres":
		return getPostgresDSN(settings, woDBName)
	case "mysql":
		return getMySQLDSN(settings, woDBName)
	}

	return ""
}

func getSQLServerDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s", settings.User, settings.Password, settings.Host, settings.Port)
	}

	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}

func getPostgresDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", settings.Host, settings.Port, settings.User, settings.Password)
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.Host, settings.Port, settings.User, settings.Password, settings.DBName)
}

func getMySQLDSN(settings *config.DBSettings, woDBName bool) string {
	if woDBName {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)", settings.User, settings.Password, settings.Host, settings.Port)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
}
