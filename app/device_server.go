package app

import (
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

// DeviceServerApplication is the application for the device server.
type DeviceServerApplication struct {
	Log    *log.Logger
	Config *config.DeviceServerSettings
	Store  *store.Store
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

	err = a.setStore()
	if err != nil {
		a.Log.Error("failed to set a store")
		return err
	}

	return nil
}

// setConfig sets the configuration for the device server application.
func (a *DeviceServerApplication) setConfig() error {
	config, err := config.ReadDeviceServerConfig()
	if err != nil {
		a.Log.Err(err, "failed to read the device server config")
		return err
	}
	a.Config = config

	return nil
}

// setStore sets the store for the device server application.
func (a *DeviceServerApplication) setStore() error {
	store, err := store.NewStore(a.Log)
	if err != nil {
		a.Log.Error("failed to create a new store")
		return err
	}

	a.Store = store

	return nil
}
