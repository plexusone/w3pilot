//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementClickTimeout time.Duration

var elementClickCmd = &cobra.Command{
	Use:   "click <selector>",
	Short: "Click an element",
	Long: `Click an element identified by CSS selector or element ref (@e1, @e2, etc.).

Examples:
  w3pilot element click "#submit"
  w3pilot element click "button.login"
  w3pilot element click "[data-testid='submit-btn']"
  w3pilot element click @e1   # Use ref from 'w3pilot map'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selectorOrRef := args[0]

		// Resolve @ref to selector if needed
		selector, err := resolveRef(selectorOrRef)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), elementClickTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		if err := el.Click(ctx, nil); err != nil {
			return fmt.Errorf("click failed: %w", err)
		}

		fmt.Printf("Clicked: %s\n", selectorOrRef)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementClickCmd)
	elementClickCmd.Flags().DurationVar(&elementClickTimeout, "timeout", 10*time.Second, "Timeout for finding and clicking element")
}
