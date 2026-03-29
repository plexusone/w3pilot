//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementDblclickTimeout time.Duration

var elementDblclickCmd = &cobra.Command{
	Use:   "dblclick <selector>",
	Short: "Double-click an element",
	Long: `Double-click an element identified by CSS selector.

Examples:
  w3pilot element dblclick "#item"
  w3pilot element dblclick ".editable"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementDblclickTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.DblClick(ctx, nil); err != nil {
			return fmt.Errorf("double-click failed: %w", err)
		}

		fmt.Printf("Double-clicked: %s\n", selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementDblclickCmd)
	elementDblclickCmd.Flags().DurationVar(&elementDblclickTimeout, "timeout", 10*time.Second, "Timeout")
}
