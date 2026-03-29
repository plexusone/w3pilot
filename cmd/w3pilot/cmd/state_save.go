package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/plexusone/w3pilot/state"
)

var stateSaveTimeout time.Duration

// StateSaveResult represents the result of saving state.
type StateSaveResult struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var stateSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current browser state",
	Long: `Save the current browser state to a named snapshot.

The state includes cookies, localStorage, and sessionStorage for all origins.
State names must be alphanumeric with dashes and underscores only.

Examples:
  w3pilot state save my-session
  w3pilot state save logged-in-user
  w3pilot state save test-data --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), stateSaveTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		// Get current state
		storageState, err := pilot.StorageState(ctx)
		if err != nil {
			return fmt.Errorf("failed to get storage state: %w", err)
		}

		// Save to file
		mgr, err := state.NewManager("")
		if err != nil {
			return fmt.Errorf("failed to create state manager: %w", err)
		}

		if err := mgr.Save(name, storageState); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}

		Output(StateSaveResult{
			Name:    name,
			Message: fmt.Sprintf("State saved as '%s'", name),
		}, func(data interface{}) string {
			return data.(StateSaveResult).Message
		})
		return nil
	},
}

func init() {
	stateCmd.AddCommand(stateSaveCmd)
	stateSaveCmd.Flags().DurationVar(&stateSaveTimeout, "timeout", 10*time.Second, "Timeout")
}
