package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/plexusone/w3pilot/state"
)

// StateListResult represents the result of listing states.
type StateListResult struct {
	States []state.StateInfo `json:"states"`
	Count  int               `json:"count"`
}

var stateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved browser states",
	Long: `List all saved browser state snapshots.

Shows the name, creation time, and size of each saved state.

Examples:
  w3pilot state list
  w3pilot state list --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := state.NewManager("")
		if err != nil {
			return fmt.Errorf("failed to create state manager: %w", err)
		}

		states, err := mgr.List()
		if err != nil {
			return fmt.Errorf("failed to list states: %w", err)
		}

		Output(StateListResult{
			States: states,
			Count:  len(states),
		}, func(data interface{}) string {
			r := data.(StateListResult)
			return formatStateList(r)
		})
		return nil
	},
}

func formatStateList(r StateListResult) string {
	if len(r.States) == 0 {
		return "No saved states"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Saved states (%d):\n\n", r.Count))

	for _, s := range r.States {
		sb.WriteString(fmt.Sprintf("  %s\n", s.Name))
		sb.WriteString(fmt.Sprintf("    Created: %s\n", s.CreatedAt.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("    Size: %d bytes\n", s.Size))
		sb.WriteString(fmt.Sprintf("    Cookies: %d\n", s.NumCookies))
		if len(s.Origins) > 0 {
			sb.WriteString(fmt.Sprintf("    Origins: %s\n", strings.Join(s.Origins, ", ")))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func init() {
	stateCmd.AddCommand(stateListCmd)
}
