// Package cmd provides the command line interface for the Doctra Middleware.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI tool.
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Doctra Middleware",
	Long:  "This is Doctra Middleware CLI tool. It is used to manage the Doctra Middleware.",
	RunE:  apiServerCommand,
}

// Init initializes the CLI tool.
func init() {
	rootCmd.AddCommand(ApiServerCommand)
	rootCmd.AddCommand(DeviceServerCommand)
}

// Execute executes the CLI tool.
func Execute() {
	cobra.CheckErr(rootCmd.ExecuteContext(context.Background()))
}

// waitForShutdown waits for the shutdown signal.
func waitForShutdown() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChan

	fmt.Println("Shutting down...")
}
