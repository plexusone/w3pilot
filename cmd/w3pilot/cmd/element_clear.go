//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementClearTimeout time.Duration

var elementClearCmd = &cobra.Command{
	Use:   "clear <selector>",
	Short: "Clear an input element",
	Long: `Clear the value of an input element.

Examples:
  w3pilot element clear "#search"
  w3pilot element clear "input[name='email']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementClearTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Clear(ctx, nil); err != nil {
			return fmt.Errorf("clear failed: %w", err)
		}

		fmt.Printf("Cleared: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementClearCmd)
	elementClearCmd.Flags().DurationVar(&elementClearTimeout, "timeout", 10*time.Second, "Timeout")
}
