package cmd

import (
	"github.com/spf13/cobra"
)

// tabCmd represents the tab command group
var tabCmd = &cobra.Command{
	Use:   "tab",
	Short: "Tab management commands",
	Long: `Commands for managing browser tabs.

Examples:
  w3pilot tab list                             # List all tabs
  w3pilot tab new                              # Create new tab
  w3pilot tab select 0                         # Switch to tab by index
  w3pilot tab close                            # Close current tab`,
}

func init() {
	rootCmd.AddCommand(tabCmd)
}
