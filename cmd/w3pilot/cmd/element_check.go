//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementCheckTimeout time.Duration

var elementCheckCmd = &cobra.Command{
	Use:   "check <selector>",
	Short: "Check a checkbox or radio button",
	Long: `Check a checkbox or radio button element.

Examples:
  w3pilot element check "#agree"
  w3pilot element check "input[type='checkbox']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementCheckTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Check(ctx, nil); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}

		fmt.Printf("Checked: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementCheckCmd)
	elementCheckCmd.Flags().DurationVar(&elementCheckTimeout, "timeout", 10*time.Second, "Timeout")
}
