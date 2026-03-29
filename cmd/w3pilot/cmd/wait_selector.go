//nolint:dupl // wait commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var waitSelectorTimeout time.Duration

var waitSelectorCmd = &cobra.Command{
	Use:   "selector <selector>",
	Short: "Wait for element to appear",
	Long: `Wait for an element matching the selector to appear in the DOM.

Examples:
  w3pilot wait selector "#modal"
  w3pilot wait selector ".loading-complete" --timeout 30s`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), waitSelectorTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		_, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		fmt.Printf("Element found: %s\n", selector)
		return nil
	},
}

func init() {
	waitCmd.AddCommand(waitSelectorCmd)
	waitSelectorCmd.Flags().DurationVar(&waitSelectorTimeout, "timeout", 30*time.Second, "Timeout")
}
