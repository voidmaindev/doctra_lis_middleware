package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewApiServer creates a new server for the API server.
func NewApiServer() (*Server, error) {
	srv, err := newServer(&app.APIServerApplication{})

	if err != nil {
		return nil, err
	}

	return srv, nil
}
