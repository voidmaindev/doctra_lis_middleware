package config

// deviceServerConfigFile is the path to the log configuration file
const deviceServerConfigFile = "config/device_server_config.json"

// DeviceServerSettings is the struct that holds the log settings
type DeviceServerSettings struct {
	Port       string
	DBSettings *DBSettings
}

// ReadDeviceServerConfig reads the log configuration file
func ReadDeviceServerConfig() (*DeviceServerSettings, error) {
	dbSettings, err := ReadDBConfig()
	if err != nil {
		return nil, err
	}

	settings := DeviceServerSettings{}
	err = ReadConfig(deviceServerConfigFile, &settings)
	if err != nil {
		return nil, err
	}

	settings.DBSettings = dbSettings

	return &settings, nil
}
