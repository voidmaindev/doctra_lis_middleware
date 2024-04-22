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
	Run:   apiServerCommand,
}

// apiServerCommand is the function that is called when the api command is executed.
func apiServerCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Starting Api Server...")

	srv := server.NewApiServer()
	srv.Start()
	defer srv.Stop()

	waitForShutdown()
}
