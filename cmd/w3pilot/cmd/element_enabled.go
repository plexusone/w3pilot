//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementEnabledTimeout time.Duration

// ElementEnabledResult represents the result of checking element enabled state.
type ElementEnabledResult struct {
	Selector string `json:"selector"`
	Enabled  bool   `json:"enabled"`
}

var elementEnabledCmd = &cobra.Command{
	Use:   "enabled <selector>",
	Short: "Check if element is enabled",
	Long: `Check if an element is enabled (not disabled).

Examples:
  w3pilot element enabled "#submit"
  w3pilot element enabled "button.next"
  w3pilot element enabled "#btn" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementEnabledTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		enabled, err := el.IsEnabled(ctx)
		if err != nil {
			return fmt.Errorf("failed to check enabled state: %w", err)
		}

		Output(ElementEnabledResult{Selector: selector, Enabled: enabled}, func(data interface{}) string {
			return fmt.Sprintf("%v", data.(ElementEnabledResult).Enabled)
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementEnabledCmd)
	elementEnabledCmd.Flags().DurationVar(&elementEnabledTimeout, "timeout", 10*time.Second, "Timeout")
}
