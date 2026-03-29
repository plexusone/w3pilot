//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageURLTimeout time.Duration

var pageURLCmd = &cobra.Command{
	Use:   "url",
	Short: "Get the current page URL",
	Long: `Get the URL of the current page.

Examples:
  w3pilot page url`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageURLTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		url, err := pilot.URL(ctx)
		if err != nil {
			return fmt.Errorf("failed to get URL: %w", err)
		}

		fmt.Println(url)
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageURLCmd)
	pageURLCmd.Flags().DurationVar(&pageURLTimeout, "timeout", 10*time.Second, "Timeout")
}
