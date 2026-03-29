//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageNewTimeout time.Duration

var pageNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new page (tab)",
	Long: `Create a new browser page/tab.

Examples:
  w3pilot page new`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageNewTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		newPage, err := pilot.NewPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to create new page: %w", err)
		}

		fmt.Printf("New page created: %s\n", newPage.BrowsingContext())
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageNewCmd)
	pageNewCmd.Flags().DurationVar(&pageNewTimeout, "timeout", 10*time.Second, "Timeout")
}
