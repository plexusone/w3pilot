//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementFocusTimeout time.Duration

var elementFocusCmd = &cobra.Command{
	Use:   "focus <selector>",
	Short: "Focus an element",
	Long: `Focus an element.

Examples:
  w3pilot element focus "#input"
  w3pilot element focus "input[name='email']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementFocusTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Focus(ctx, nil); err != nil {
			return fmt.Errorf("focus failed: %w", err)
		}

		fmt.Printf("Focused: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementFocusCmd)
	elementFocusCmd.Flags().DurationVar(&elementFocusTimeout, "timeout", 10*time.Second, "Timeout")
}
