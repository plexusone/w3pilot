//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageBackTimeout time.Duration

var pageBackCmd = &cobra.Command{
	Use:   "back",
	Short: "Navigate back in browser history",
	Long: `Navigate back to the previous page in browser history.

Examples:
  w3pilot page back`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageBackTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.Back(ctx); err != nil {
			return fmt.Errorf("back navigation failed: %w", err)
		}

		url, _ := pilot.URL(ctx)
		fmt.Println("Navigated back")
		if url != "" {
			fmt.Printf("URL: %s\n", url)
		}

		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageBackCmd)
	pageBackCmd.Flags().DurationVar(&pageBackTimeout, "timeout", 10*time.Second, "Timeout")
}
