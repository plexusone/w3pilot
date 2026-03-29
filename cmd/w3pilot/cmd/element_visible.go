//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementVisibleTimeout time.Duration

var elementVisibleCmd = &cobra.Command{
	Use:   "visible <selector>",
	Short: "Check if element is visible",
	Long: `Check if an element is visible on the page.

Examples:
  w3pilot element visible "#modal"
  w3pilot element visible ".loading-spinner"`,
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

		fmt.Println(visible)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementVisibleCmd)
	elementVisibleCmd.Flags().DurationVar(&elementVisibleTimeout, "timeout", 10*time.Second, "Timeout")
}
