//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageReloadTimeout time.Duration

var pageReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the current page",
	Long: `Reload the current page.

Examples:
  w3pilot page reload`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageReloadTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.Reload(ctx); err != nil {
			return fmt.Errorf("reload failed: %w", err)
		}

		fmt.Println("Page reloaded")
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageReloadCmd)
	pageReloadCmd.Flags().DurationVar(&pageReloadTimeout, "timeout", 30*time.Second, "Timeout")
}
