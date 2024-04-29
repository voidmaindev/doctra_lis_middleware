package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidmaindev/doctra_lis_middleware/server"
)

// DeviceServerCommand is the command to start the Doctra Middleware Device server.
var DeviceServerCommand = &cobra.Command{
	Use:   "device",
	Short: "Start the Doctra Middleware Device server",
	Long:  "This command starts the Doctra Middleware Device server.",
	RunE:  deviceServerCommand,
}

// deviceServerCommand is the function that is called when the device command is executed.
func deviceServerCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("creating Device Server...")

	srv, err := server.NewDeviceServer()
	if err != nil {
		return err
	}

	err = srv.Start()
	if err != nil {
		srv.Log.Fatal("failed to start the Device server")
	}

	defer srv.Stop()

	waitForShutdown()

	return nil
}
