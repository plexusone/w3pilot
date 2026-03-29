//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementValueTimeout time.Duration

// ElementValueResult represents the result of getting element value.
type ElementValueResult struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

var elementValueCmd = &cobra.Command{
	Use:   "value <selector>",
	Short: "Get input element value",
	Long: `Get the value of an input element.

Examples:
  w3pilot element value "#email"
  w3pilot element value "input[name='search']"
  w3pilot element value "#search" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementValueTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		value, err := el.Value(ctx)
		if err != nil {
			return fmt.Errorf("failed to get value: %w", err)
		}

		Output(ElementValueResult{Selector: selector, Value: value}, func(data interface{}) string {
			return data.(ElementValueResult).Value
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementValueCmd)
	elementValueCmd.Flags().DurationVar(&elementValueTimeout, "timeout", 10*time.Second, "Timeout")
}
