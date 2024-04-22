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
	Run:   deviceServerCommand,
}

// deviceServerCommand is the function that is called when the device command is executed.
func deviceServerCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Starting Device Server...")

	srv := server.NewDeviceServer()
	srv.Start()
	defer srv.Stop()

	waitForShutdown()
}
