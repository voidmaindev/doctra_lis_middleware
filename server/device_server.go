package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewDeviceServer creates a new server for the device server.
func NewDeviceServer() (*Server, error) {
	srv, err := newServer(&app.DeviceServerApplication{})

	if err != nil {
		return nil, err
	}

	return srv, nil
}
