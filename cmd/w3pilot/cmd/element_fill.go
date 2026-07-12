//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementFillTimeout time.Duration

var elementFillCmd = &cobra.Command{
	Use:   "fill <selector> <value>",
	Short: "Fill an input with a value",
	Long: `Fill an input element with a value.
This clears the input first and then sets the value directly.
Faster than type but doesn't simulate key events.

Examples:
  w3pilot element fill "#password" "secret123"
  w3pilot element fill "textarea" "Long text content..."
  w3pilot element fill @e3 "john@example.com"   # Use ref from 'w3pilot map'`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selectorOrRef := args[0]
		value := args[1]

		// Resolve @ref to selector if needed
		selector, err := resolveRef(selectorOrRef)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), elementFillTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Fill(ctx, value, nil); err != nil {
			return fmt.Errorf("fill failed: %w", err)
		}

		fmt.Printf("Filled: %s\n", selectorOrRef)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementFillCmd)
	elementFillCmd.Flags().DurationVar(&elementFillTimeout, "timeout", 10*time.Second, "Timeout")
}
