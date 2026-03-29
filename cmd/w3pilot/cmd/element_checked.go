//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementCheckedTimeout time.Duration

var elementCheckedCmd = &cobra.Command{
	Use:   "checked <selector>",
	Short: "Check if checkbox/radio is checked",
	Long: `Check if a checkbox or radio button is checked.

Examples:
  w3pilot element checked "#agree"
  w3pilot element checked "input[type='checkbox']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementCheckedTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		checked, err := el.IsChecked(ctx)
		if err != nil {
			return fmt.Errorf("failed to check state: %w", err)
		}

		fmt.Println(checked)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementCheckedCmd)
	elementCheckedCmd.Flags().DurationVar(&elementCheckedTimeout, "timeout", 10*time.Second, "Timeout")
}
