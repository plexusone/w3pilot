package cmd

import (
	"github.com/spf13/cobra"
)

// stateCmd represents the state command group
var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Manage browser state snapshots",
	Long: `Commands for saving, loading, and managing browser state snapshots.

State snapshots preserve cookies, localStorage, and sessionStorage,
allowing you to quickly restore authenticated sessions or test states.

Examples:
  w3pilot state save my-session            # Save current state
  w3pilot state list                       # List saved states
  w3pilot state load my-session            # Load a saved state
  w3pilot state delete my-session          # Delete a saved state`,
}

func init() {
	rootCmd.AddCommand(stateCmd)
}
