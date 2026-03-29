//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementPressTimeout time.Duration

var elementPressCmd = &cobra.Command{
	Use:   "press <selector> <key>",
	Short: "Press a key on an element",
	Long: `Press a key on a focused element.

Key can be: Enter, Tab, Escape, ArrowDown, ArrowUp, etc.

Examples:
  w3pilot element press "#input" Enter
  w3pilot element press "#search" Tab`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		key := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), elementPressTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Press(ctx, key, nil); err != nil {
			return fmt.Errorf("press failed: %w", err)
		}

		fmt.Printf("Pressed %s on: %s\n", key, selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementPressCmd)
	elementPressCmd.Flags().DurationVar(&elementPressTimeout, "timeout", 10*time.Second, "Timeout")
}
