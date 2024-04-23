package app

import (
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// ApiServerApplication is the application for the API server.
type ApiServerApplication struct {
	Log    *log.Logger
	Config *config.ApiServerSettings
}

// SetLogger sets the logger for the API server application.
func (a *ApiServerApplication) SetLogger(l *log.Logger) {
	a.Log = l
}

// InitApp initializes the API server application.
func (a *ApiServerApplication) InitApp() error {
	config, err := config.ReadApiServerConfig()
	if err != nil {
		return err
	}
	a.Config = config

	return nil
}
