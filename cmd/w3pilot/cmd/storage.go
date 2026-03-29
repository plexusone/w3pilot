package cmd

import (
	"github.com/spf13/cobra"
)

// storageCmd represents the storage command group
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Browser storage commands",
	Long: `Commands for managing cookies, localStorage, and sessionStorage.

Examples:
  w3pilot storage cookies-get                  # Get all cookies
  w3pilot storage cookies-set name=value       # Set a cookie
  w3pilot storage local-get key                # Get localStorage item
  w3pilot storage local-set key value          # Set localStorage item
  w3pilot storage session-get key              # Get sessionStorage item
  w3pilot storage state-get state.json         # Save storage state
  w3pilot storage clear-all                    # Clear all storage`,
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
