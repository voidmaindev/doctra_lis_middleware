package server

import (
	"github.com/voidmaindev/doctra_lis_middleware/app"
)

// server is the struct that represents the server.
type server struct {
	App    app.App
	Log    int
	Config int
}

// NewApiServer creates a new server for the API server.
func newServer(a app.App) *server {
	//implement logger

	//implement config

	return &server{App: a}
}

// Start starts the server.
func (s *server) Start() error {
	err := s.App.InitApp()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the server.
func (s *server) Stop() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()
	// s.Shutdown(ctx)

	// log.Println("Shutting down")
	// os.Exit(0)
}
