//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementUncheckTimeout time.Duration

var elementUncheckCmd = &cobra.Command{
	Use:   "uncheck <selector>",
	Short: "Uncheck a checkbox",
	Long: `Uncheck a checkbox element.

Examples:
  w3pilot element uncheck "#newsletter"
  w3pilot element uncheck "input[name='subscribe']"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementUncheckTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Uncheck(ctx, nil); err != nil {
			return fmt.Errorf("uncheck failed: %w", err)
		}

		fmt.Printf("Unchecked: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementUncheckCmd)
	elementUncheckCmd.Flags().DurationVar(&elementUncheckTimeout, "timeout", 10*time.Second, "Timeout")
}
