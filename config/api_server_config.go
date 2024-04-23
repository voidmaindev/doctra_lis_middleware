package config

// apiServerConfigFile is the path to the log configuration file
const apiServerConfigFile = "config/api_server_config.json"

// ApiServerSettings is the struct that holds the log settings
type ApiServerSettings struct {
	APIPort    string
	DBSettings *DBSettings
}

// ReadApiServerConfig reads the log configuration file
func ReadApiServerConfig() (*ApiServerSettings, error) {
	dbSettings, err := ReadDBConfig()
	if err != nil {
		return nil, err
	}

	settings := ApiServerSettings{}
	err = ReadConfig(apiServerConfigFile, &settings)
	if err != nil {
		return nil, err
	}

	settings.DBSettings = dbSettings

	return &settings, nil
}
