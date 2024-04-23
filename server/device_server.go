package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewDeviceServer creates a new server for the device server.
func NewDeviceServer() *Server {
	srv, err := newServer(&app.DeviceServerApplication{})

	if err != nil {
		panic(err)
	}

	return srv
}
