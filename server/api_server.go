package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewApiServer creates a new server for the API server.
func NewApiServer() *Server {
	srv, err := newServer(&app.ApiServerApplication{})

	if err != nil {
		panic(err)
	}

	return srv
}
