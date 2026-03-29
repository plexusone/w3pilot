//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageContentTimeout time.Duration

// ContentResult represents the result of getting the page content.
type ContentResult struct {
	Content string `json:"content"`
}

var pageContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Get the page HTML content",
	Long: `Get the full HTML content of the current page.

Examples:
  w3pilot page content
  w3pilot page content > page.html
  w3pilot page content --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageContentTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		content, err := pilot.Content(ctx)
		if err != nil {
			return fmt.Errorf("failed to get content: %w", err)
		}

		Output(ContentResult{Content: content}, func(data interface{}) string {
			return data.(ContentResult).Content
		})
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageContentCmd)
	pageContentCmd.Flags().DurationVar(&pageContentTimeout, "timeout", 30*time.Second, "Timeout")
}
