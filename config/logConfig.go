package config

// logConfigFile is the path to the log configuration file
const logConfigFile = "config/logConfig.json"

// LogSettings is the struct that holds the log settings
type LogSettings struct {
	Disable    bool
	Output     string
	TimeFormat string
	AddCaller  bool
	AddPid     bool
}

// ReadLogConfig reads the log configuration file
func ReadLogConfig() (*LogSettings, error) {
	settings := LogSettings{}
	err := ReadConfig(logConfigFile, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
