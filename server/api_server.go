package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewApiServer creates a new server for the API server.
func NewApiServer() *server {
	srv, err := newServer(&app.ApiServerApplication{})

	if err != nil {
		panic(err)
	}

	return srv
}
