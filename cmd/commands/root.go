package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI tool.
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Doctra Middleware",
	Long:  "This is Doctra Middleware CLI tool. It is used to manage the Doctra Middleware.",
}

// Init initializes the CLI tool.
func init() {
	rootCmd.AddCommand(ApiServerCommand)
}

// Execute executes the CLI tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
