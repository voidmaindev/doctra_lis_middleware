// Package server provides the server struct that represents the server.
package server

import (
	"github.com/voidmaindev/doctra_lis_middleware/app"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// server is the struct that represents the server.
type Server struct {
	App app.App
	Log *log.Logger
}

// NewServer creates a new server.
func newServer(a app.App) (*Server, error) {
	logger, err := log.NewLogger()
	if err != nil {
		return nil, err
	}

	srv := &Server{
		App: a,
		Log: logger,
	}

	return srv, nil
}

// Start starts the server.
func (s *Server) Start() error {
	s.Log.Info("Starting server...")

	s.App.SetLogger(s.Log)

	err := s.App.InitApp()
	if err != nil {
		s.Log.Error("failed to initialize the application")
		return err
	}

	err = s.App.Start()
	if err != nil {
		s.Log.Error("failed to start the application")
		return err
	}

	return nil
}

// Stop stops the server.
func (s *Server) Stop() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()
	// s.Shutdown(ctx)

	// log.Println("Shutting down")
	// os.Exit(0)
}
