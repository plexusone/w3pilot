package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/plexusone/w3pilot/state"
)

var stateLoadTimeout time.Duration

// StateLoadResult represents the result of loading state.
type StateLoadResult struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var stateLoadCmd = &cobra.Command{
	Use:   "load <name>",
	Short: "Load a saved browser state",
	Long: `Load a previously saved browser state by name.

This restores cookies, localStorage, and sessionStorage.
You may need to reload the page for changes to take effect.

Examples:
  w3pilot state load my-session
  w3pilot state load logged-in-user
  w3pilot state load test-data --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), stateLoadTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		// Load from file
		mgr, err := state.NewManager("")
		if err != nil {
			return fmt.Errorf("failed to create state manager: %w", err)
		}

		storageState, err := mgr.Load(name)
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}

		// Apply state
		if err := pilot.SetStorageState(ctx, storageState); err != nil {
			return fmt.Errorf("failed to apply storage state: %w", err)
		}

		Output(StateLoadResult{
			Name:    name,
			Message: fmt.Sprintf("State '%s' loaded", name),
		}, func(data interface{}) string {
			return data.(StateLoadResult).Message
		})
		return nil
	},
}

func init() {
	stateCmd.AddCommand(stateLoadCmd)
	stateLoadCmd.Flags().DurationVar(&stateLoadTimeout, "timeout", 10*time.Second, "Timeout")
}
