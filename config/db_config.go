package config

// dbConfigFile is the path to the log configuration file
const dbConfigFile = "config/log_config.json"

// LogSettings is the struct that holds the log settings
type DBSettings struct {
	DriverName string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

// ReadLogConfig reads the log configuration file
func ReadDBConfig() (*DBSettings, error) {
	settings := DBSettings{}
	err := ReadConfig(dbConfigFile, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
