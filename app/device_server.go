package app

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
	"github.com/voidmaindev/doctra_lis_middleware/tcp"
)

// DeviceServerApplication is the application for the device server.
type DeviceServerApplication struct {
	Log      *log.Logger
	Config   *config.DeviceServerSettings
	Listener net.Listener
	Store    *store.Store
	TCP      *tcp.TCP
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

	err = a.setListener()
	if err != nil {
		a.Log.Error("failed to set the TCP listener")
		return err
	}

	err = a.setStore()
	if err != nil {
		a.Log.Error("failed to set a store")
		return err
	}

	err = a.setTCP()
	if err != nil {
		a.Log.Error("failed to set the TCP")
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

// setListener sets the listener for the device server application.
func (a *DeviceServerApplication) setListener() error {
	listener, err := net.Listen("tcp", a.Config.Port)
	if err != nil {
		a.Log.Err(err, "failed to start the TCP listener")
		return err
	}

	a.Listener = listener

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

// setTCP sets the TCP for the device server application.
func (a *DeviceServerApplication) setTCP() error {
	tcp := tcp.NewTCP(a.Log, a.Listener)
	a.TCP = tcp

	return nil
}

// Run runs the device server application.
func (a *DeviceServerApplication) Start() error {
	a.Log.Info("starting the device server")

	go a.TCP.Accept()

	for msg := range a.TCP.RcvChannel {
		if a.TCP.Conns[msg.ConnString] == nil {
			continue
		}

		a.Log.Info("received a message from " + msg.ConnString)
		a.Log.Info(string(msg.Data))
	}

	return nil
}

// Stop stops the device server application.
func (a *DeviceServerApplication) Stop() error {
	a.Log.Info("stopping the device server")

	err := a.Listener.Close()
	if err != nil {
		a.Log.Err(err, "failed to stop the TCP listener")
		return err
	}

	close(a.TCP.RcvChannel)

	return nil
}
