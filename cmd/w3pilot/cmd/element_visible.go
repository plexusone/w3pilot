//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementVisibleTimeout time.Duration

// ElementVisibleResult represents the result of checking element visibility.
type ElementVisibleResult struct {
	Selector string `json:"selector"`
	Visible  bool   `json:"visible"`
}

var elementVisibleCmd = &cobra.Command{
	Use:   "visible <selector>",
	Short: "Check if element is visible",
	Long: `Check if an element is visible on the page.

Examples:
  w3pilot element visible "#modal"
  w3pilot element visible ".loading-spinner"
  w3pilot element visible "#dialog" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementVisibleTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		visible, err := el.IsVisible(ctx)
		if err != nil {
			return fmt.Errorf("failed to check visibility: %w", err)
		}

		Output(ElementVisibleResult{Selector: selector, Visible: visible}, func(data interface{}) string {
			return fmt.Sprintf("%v", data.(ElementVisibleResult).Visible)
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementVisibleCmd)
	elementVisibleCmd.Flags().DurationVar(&elementVisibleTimeout, "timeout", 10*time.Second, "Timeout")
}
