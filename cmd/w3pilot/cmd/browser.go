package cmd

import (
	"github.com/spf13/cobra"
)

// browserCmd represents the browser command group
var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "Browser lifecycle commands",
	Long: `Commands for managing browser instances.

Examples:
  w3pilot browser launch              # Launch visible browser
  w3pilot browser launch --headless   # Launch headless browser
  w3pilot browser quit                # Close browser`,
}

func init() {
	rootCmd.AddCommand(browserCmd)
}
