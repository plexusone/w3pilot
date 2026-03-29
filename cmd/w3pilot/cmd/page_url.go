//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageURLTimeout time.Duration

// URLResult represents the result of getting the page URL.
type URLResult struct {
	URL string `json:"url"`
}

var pageURLCmd = &cobra.Command{
	Use:   "url",
	Short: "Get the current page URL",
	Long: `Get the URL of the current page.

Examples:
  w3pilot page url
  w3pilot page url --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageURLTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		url, err := pilot.URL(ctx)
		if err != nil {
			return fmt.Errorf("failed to get URL: %w", err)
		}

		Output(URLResult{URL: url}, func(data interface{}) string {
			return data.(URLResult).URL
		})
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageURLCmd)
	pageURLCmd.Flags().DurationVar(&pageURLTimeout, "timeout", 10*time.Second, "Timeout")
}
