//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementHoverTimeout time.Duration

var elementHoverCmd = &cobra.Command{
	Use:   "hover <selector>",
	Short: "Hover over an element",
	Long: `Hover the mouse over an element.

Examples:
  w3pilot element hover "#menu"
  w3pilot element hover ".dropdown-trigger"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementHoverTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Hover(ctx, nil); err != nil {
			return fmt.Errorf("hover failed: %w", err)
		}

		fmt.Printf("Hovering: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementHoverCmd)
	elementHoverCmd.Flags().DurationVar(&elementHoverTimeout, "timeout", 10*time.Second, "Timeout")
}
