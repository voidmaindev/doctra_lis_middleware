package app

import (
	"fmt"
	"net"
	"strings"

	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/driver"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/services"
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
	address := a.Config.Host + ":" + a.Config.Port
	listener, err := net.Listen("tcp", address)
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

	go a.TCP.AcceptConnections()
	go a.ManageMessages()

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

// ManageMessages manages the messages received by the device server.
func (a *DeviceServerApplication) ManageMessages() {
	for msg := range a.TCP.RcvChannel {
		conn := a.TCP.Conns[msg.ConnString]
		if conn == nil {
			a.Log.Error("failed to get a connection by network address: " + msg.ConnString)
			continue
		}

		a.Log.Info("received a message from " + msg.ConnString)
		a.Log.Info(string(msg.Data))

		device, err := a.Store.DeviceStore.GetByNetAddress(msg.ConnString)
		if err != nil {
			a.Log.Error("failed to get a device by network address: " + msg.ConnString)
			continue
		}

		err = a.processDeviceMessage(msg.Data, conn, device, a.Log, a.Store)
		if err != nil {
			a.Log.Error("failed to process the device message")
			continue
		}
	}
}

// processDeviceMessage processes the device message.
func (a *DeviceServerApplication) processDeviceMessage(deviceMsg []byte, conn *tcp.ConnData, device *model.Device, log *log.Logger, store *store.Store) error {
	device, err := store.DeviceStore.GetByID(device.ID)
	if err != nil {
		log.Error("failed to get a device by ID: " + fmt.Sprint(device.ID))
		return err
	}

	deviceQueryService := services.NewDeviceQueryService(a.Config.DBSettings.QueryHost, device.Serial)
	
	deviceDriver, err := driver.NewDriver(device.DeviceModel.Driver, log, store, deviceQueryService)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create a driver for %s with driver name %s", device.Name, device.DeviceModel.Driver))
		return err
	}

	err = deviceDriver.SendSimpleACK(conn.Conn)
	if err != nil {
		log.Err(err, "failed to send an ACK message to "+device.Name)
		return err
	}

	msg := string(deviceMsg)
	for k, v := range deviceDriver.DataToBeReplaced() {
		msg = strings.ReplaceAll(msg, k, v)
	}

	rawDatas := driver.GetRawDatas(deviceDriver, msg, conn.PrevData)

	for _, rawData := range rawDatas {
		rd := &model.RawData{
			ConnString: conn.ConnString,
			DeviceID:   device.ID,
			Data:       []byte(rawData),
			Processed:  true,
		}

		labDatas, additionalData, err := deviceDriver.Unmarshal(rawData)
		if err != nil {
			deviceDriver.Log().Error("failed to unmarshal a raw data from " + device.Name)
			rd.Processed = false
		}

		err = deviceDriver.Store().RawDataStore.Create(rd)
		if err != nil {
			deviceDriver.Log().Error("failed to create a raw data from " + device.Name)
			return err
		}

		for _, labData := range labDatas {
			labData.RawDataID = rd.ID
			labData.DeviceID = device.ID

			err = deviceDriver.Store().LabDataStore.Create(labData)
			if err != nil {
				deviceDriver.Log().Error(fmt.Sprintf("failed to create a lab data from %s with barcode %s and index %d", device.Name, labData.Barcode, labData.Index))
				rd.Processed = false
				continue
			}

			if rd.Processed {
				deviceDriver.PostUnmarshalActions(conn.Conn, additionalData)
			}
		}
	}

	return nil
}
