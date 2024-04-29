package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidmaindev/doctra_lis_middleware/server"
)

// ApiServerCommand is the command to start the Doctra Middleware API server.
var ApiServerCommand = &cobra.Command{
	Use:   "api",
	Short: "Start the Doctra Middleware API server",
	Long:  "This command starts the Doctra Middleware API server.",
	RunE:  apiServerCommand,
}

// apiServerCommand is the function that is called when the api command is executed.
func apiServerCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("creating Api Server...")

	srv, err := server.NewApiServer()
	if err != nil {
		return err
	}

	err = srv.Start()
	if err != nil {
		srv.Log.Fatal("failed to start the API server")
	}

	defer srv.Stop()

	waitForShutdown()

	return nil
}
