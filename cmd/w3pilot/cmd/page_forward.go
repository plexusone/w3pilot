//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageForwardTimeout time.Duration

var pageForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Navigate forward in browser history",
	Long: `Navigate forward to the next page in browser history.

Examples:
  w3pilot page forward`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageForwardTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.Forward(ctx); err != nil {
			return fmt.Errorf("forward navigation failed: %w", err)
		}

		url, _ := pilot.URL(ctx)
		fmt.Println("Navigated forward")
		if url != "" {
			fmt.Printf("URL: %s\n", url)
		}

		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageForwardCmd)
	pageForwardCmd.Flags().DurationVar(&pageForwardTimeout, "timeout", 10*time.Second, "Timeout")
}
