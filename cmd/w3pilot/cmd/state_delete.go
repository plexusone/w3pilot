package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/plexusone/w3pilot/state"
)

// StateDeleteResult represents the result of deleting state.
type StateDeleteResult struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var stateDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a saved browser state",
	Long: `Delete a previously saved browser state by name.

Examples:
  w3pilot state delete my-session
  w3pilot state delete old-test-data`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		mgr, err := state.NewManager("")
		if err != nil {
			return fmt.Errorf("failed to create state manager: %w", err)
		}

		if err := mgr.Delete(name); err != nil {
			return fmt.Errorf("failed to delete state: %w", err)
		}

		Output(StateDeleteResult{
			Name:    name,
			Message: fmt.Sprintf("State '%s' deleted", name),
		}, func(data interface{}) string {
			return data.(StateDeleteResult).Message
		})
		return nil
	},
}

func init() {
	stateCmd.AddCommand(stateDeleteCmd)
}
