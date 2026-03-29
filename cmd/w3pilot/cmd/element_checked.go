//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementCheckedTimeout time.Duration

// ElementCheckedResult represents the result of checking checkbox/radio state.
type ElementCheckedResult struct {
	Selector string `json:"selector"`
	Checked  bool   `json:"checked"`
}

var elementCheckedCmd = &cobra.Command{
	Use:   "checked <selector>",
	Short: "Check if checkbox/radio is checked",
	Long: `Check if a checkbox or radio button is checked.

Examples:
  w3pilot element checked "#agree"
  w3pilot element checked "input[type='checkbox']"
  w3pilot element checked "#terms" --format json`,
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

		Output(ElementCheckedResult{Selector: selector, Checked: checked}, func(data interface{}) string {
			return fmt.Sprintf("%v", data.(ElementCheckedResult).Checked)
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementCheckedCmd)
	elementCheckedCmd.Flags().DurationVar(&elementCheckedTimeout, "timeout", 10*time.Second, "Timeout")
}
