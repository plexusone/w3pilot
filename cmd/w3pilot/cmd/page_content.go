//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageContentTimeout time.Duration

var pageContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Get the page HTML content",
	Long: `Get the full HTML content of the current page.

Examples:
  w3pilot page content
  w3pilot page content > page.html`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageContentTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		content, err := pilot.Content(ctx)
		if err != nil {
			return fmt.Errorf("failed to get content: %w", err)
		}

		fmt.Println(content)
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageContentCmd)
	pageContentCmd.Flags().DurationVar(&pageContentTimeout, "timeout", 30*time.Second, "Timeout")
}
