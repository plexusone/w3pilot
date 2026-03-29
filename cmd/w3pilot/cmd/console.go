package cmd

import (
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command group
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Console message commands",
	Long: `Commands for capturing and viewing console messages.

Examples:
  w3pilot console messages                     # Get console messages
  w3pilot console messages --level error       # Get only errors
  w3pilot console clear                        # Clear console buffer`,
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
