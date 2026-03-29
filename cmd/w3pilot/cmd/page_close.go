//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageCloseTimeout time.Duration

var pageCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close the current page",
	Long: `Close the current page (but not the browser).

Examples:
  w3pilot page close`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageCloseTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.Close(ctx); err != nil {
			return fmt.Errorf("failed to close page: %w", err)
		}

		fmt.Println("Page closed")
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageCloseCmd)
	pageCloseCmd.Flags().DurationVar(&pageCloseTimeout, "timeout", 10*time.Second, "Timeout")
}
