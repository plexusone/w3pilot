//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageTitleTimeout time.Duration

var pageTitleCmd = &cobra.Command{
	Use:   "title",
	Short: "Get the page title",
	Long: `Get the title of the current page.

Examples:
  w3pilot page title`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageTitleTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		title, err := pilot.Title(ctx)
		if err != nil {
			return fmt.Errorf("failed to get title: %w", err)
		}

		fmt.Println(title)
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageTitleCmd)
	pageTitleCmd.Flags().DurationVar(&pageTitleTimeout, "timeout", 10*time.Second, "Timeout")
}
