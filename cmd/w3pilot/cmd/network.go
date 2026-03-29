package cmd

import (
	"github.com/spf13/cobra"
)

// networkCmd represents the network command group
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network interception commands",
	Long: `Commands for intercepting and monitoring network requests.

Examples:
  w3pilot network requests                     # List captured requests
  w3pilot network route "**/*.png" --status 404  # Mock route
  w3pilot network unroute "**/*.png"           # Remove route
  w3pilot network offline true                 # Enable offline mode`,
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
