package config

const logConfigFile = "config/logConfig.json"

type LogSettings struct {
	Disable    bool
	Output     string
	TimeFormat string
	AddCaller  bool
	AddPid     bool
}

func ReadLogConfig() (*LogSettings, error) {
	settings := LogSettings{}
	err := ReadConfig(logConfigFile, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
