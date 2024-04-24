package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/api"
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// ApiServerApplication is the application for the API server.
type ApiServerApplication struct {
	Log    *log.Logger
	Config *config.ApiServerSettings
	Router *fiber.App
}

// SetLogger sets the logger for the API server application.
func (a *ApiServerApplication) SetLogger(l *log.Logger) {
	a.Log = l
}

// InitApp initializes the API server application.
func (a *ApiServerApplication) InitApp() error {
	err := a.setConfig()
	if err != nil {
		a.Log.Error("Failed to set the API server config")
		return err
	}

	a.setRouter()

	return nil
}

func (a *ApiServerApplication) setConfig() error {
	config, err := config.ReadApiServerConfig()
	if err != nil {
		a.Log.Err(err, "Failed to read the API server config")
		return err
	}
	a.Config = config

	return nil
}

func (a *ApiServerApplication) setRouter() {
	a.Router = api.NewRouter(a.Log.Logger)
}
