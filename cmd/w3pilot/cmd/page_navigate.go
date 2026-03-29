//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageNavigateTimeout time.Duration

var pageNavigateCmd = &cobra.Command{
	Use:   "navigate <url>",
	Short: "Navigate to a URL",
	Long: `Navigate the browser to the specified URL.

Examples:
  w3pilot page navigate https://example.com
  w3pilot page navigate https://google.com --timeout 30s`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), pageNavigateTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.Go(ctx, url); err != nil {
			return fmt.Errorf("navigation failed: %w", err)
		}

		title, _ := pilot.Title(ctx)
		fmt.Printf("Navigated to: %s\n", url)
		if title != "" {
			fmt.Printf("Title: %s\n", title)
		}

		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageNavigateCmd)
	pageNavigateCmd.Flags().DurationVar(&pageNavigateTimeout, "timeout", 30*time.Second, "Navigation timeout")
}
