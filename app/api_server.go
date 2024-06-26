package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/api"
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

// APIServerApplication is the application for the API server.
type APIServerApplication struct {
	Log    *log.Logger
	Config *config.ApiServerSettings
	Router *fiber.App
	Store  *store.Store
	API    *api.API
}

// SetLogger sets the logger for the API server application.
func (a *APIServerApplication) SetLogger(l *log.Logger) {
	a.Log = l
}

// InitApp initializes the API server application.
func (a *APIServerApplication) InitApp() error {
	err := a.setConfig()
	if err != nil {
		a.Log.Error("failed to set the API server config")
		return err
	}

	a.setRouter()

	err = a.setStore()
	if err != nil {
		a.Log.Error("failed to set a store")
		return err
	}

	err = a.CreateDefAdmin()
	if err != nil {
		a.Log.Error("failed to create the default admin user")
		return err
	}

	err = a.setAPI()
	if err != nil {
		a.Log.Error("failed to set the API")
		return err
	}

	return nil
}

// setConfig sets the configuration for the API server application.
func (a *APIServerApplication) setConfig() error {
	config, err := config.ReadApiServerConfig()
	if err != nil {
		a.Log.Err(err, "failed to read the API server config")
		return err
	}
	a.Config = config

	return nil
}

// setRouter sets the router for the API server application.
func (a *APIServerApplication) setRouter() {
	a.Router = api.NewRouter(a.Log.Logger)
}

// setStore sets the store for the API server application.
func (a *APIServerApplication) setStore() error {
	store, err := store.NewStore(a.Log)
	if err != nil {
		a.Log.Error("failed to create a new store")
		return err
	}

	a.Store = store

	return nil
}

// setAPI sets the API for the API server application.
func (a *APIServerApplication) setAPI() error {
	api, err := api.NewAPI(a.Log, a.Router, a.Store)
	if err != nil {
		a.Log.Error("failed to create a new API")
		return err
	}

	a.API = api

	return nil
}

// CreateDefAdmin creates the default admin user.
func (a *APIServerApplication) CreateDefAdmin() error {
	if !a.Config.DBSettings.CreateDefAdmin {
		return nil
	}

	admin, err := a.Store.UserStore.GetByUsername("admin")
	if err == nil && admin != nil {
		return nil
	}

	a.Log.Info("creating the default admin user")

	authUser := &model.AuthUser{
		Username: "admin",
		Password: "Ab123456",
		Role:     "admin",
	}

	user, err := model.NewUserFromAuthUser(authUser)
	if err != nil {
		a.Log.Err(err, "failed to create a new user from the auth user")
		return err
	}

	err = a.Store.UserStore.Create(user)
	if err != nil {
		a.Log.Error("failed to create the default admin user")
		return err
	}

	return nil
}

// Start starts the API server application.
func (a *APIServerApplication) Start() error {
	a.Log.Info("starting the API server")

	address := a.Config.Host + ":" + a.Config.Port

	go func() {
		err := a.Router.Listen(address)
		if err != nil {
			a.Log.Err(err, "failed to start the API server")
		}
	}()

	return nil
}

// Stop stops the API server application.
func (a *APIServerApplication) Stop() error {
	a.Log.Info("stopping the API server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*serverStopTimeout)
	defer cancel()
	err := a.Router.ShutdownWithContext(ctx)
	if err != nil {
		a.Log.Err(err, "failed to stop the API server")
		return err
	}

	return nil
}
