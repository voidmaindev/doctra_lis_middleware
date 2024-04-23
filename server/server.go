// Package server provides the server struct that represents the server.
package server

import (
	"github.com/voidmaindev/doctra_lis_middleware/app"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// server is the struct that represents the server.
type server struct {
	App    app.App
	Log    *log.Logger
	Config int
}

// NewServer creates a new server.
func newServer(a app.App) (*server, error) {
	logger, err := log.NewLogger()
	if err != nil {
		return nil, err
	}

	srv := &server{
		App: a,
		Log: logger,
	}

	return srv, nil
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
