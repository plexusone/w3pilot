//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementTypeTimeout time.Duration

var elementTypeCmd = &cobra.Command{
	Use:   "type <selector> <text>",
	Short: "Type text into an element",
	Long: `Type text into an input element character by character.
This simulates real keyboard input with key events.

Examples:
  w3pilot element type "#email" "test@example.com"
  w3pilot element type "input[name='search']" "query"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		text := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), elementTypeTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Type(ctx, text, nil); err != nil {
			return fmt.Errorf("type failed: %w", err)
		}

		fmt.Printf("Typed into: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementTypeCmd)
	elementTypeCmd.Flags().DurationVar(&elementTypeTimeout, "timeout", 10*time.Second, "Timeout")
}
