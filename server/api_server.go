package server

import "github.com/voidmaindev/doctra_lis_middleware/app"

// NewApiServer creates a new server for the API server.
func NewApiServer() *server {
	return (newServer(&app.ApiServerApplication{}))
}
