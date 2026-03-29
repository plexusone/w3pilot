//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementScrollTimeout time.Duration

var elementScrollCmd = &cobra.Command{
	Use:   "scroll <selector>",
	Short: "Scroll element into view",
	Long: `Scroll an element into the visible viewport.

Examples:
  w3pilot element scroll "#footer"
  w3pilot element scroll ".lazy-load"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementScrollTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.ScrollIntoView(ctx, nil); err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}

		fmt.Printf("Scrolled into view: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementScrollCmd)
	elementScrollCmd.Flags().DurationVar(&elementScrollTimeout, "timeout", 10*time.Second, "Timeout")
}
