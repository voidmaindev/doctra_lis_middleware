package app

import (
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// DeviceServerApplication is the application for the device server.
type DeviceServerApplication struct {
	Log    *log.Logger
	Config *config.DeviceServerSettings
}

// SetLogger sets the logger for the device server application.
func (a *DeviceServerApplication) SetLogger(l *log.Logger) {
	a.Log = l
}

// InitApp initializes the device server application.
func (a *DeviceServerApplication) InitApp() error {
	err := a.setConfig()
	if err != nil {
		a.Log.Error("failed to set the device server config")
		return err
	}

	return nil
}

func (a *DeviceServerApplication) setConfig() error {
	config, err := config.ReadDeviceServerConfig()
	if err != nil {
		a.Log.Err(err, "failed to read the device server config")
		return err
	}
	a.Config = config

	return nil
}
