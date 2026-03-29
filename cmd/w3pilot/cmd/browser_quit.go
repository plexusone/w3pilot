//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var browserQuitTimeout time.Duration

var browserQuitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit the browser",
	Long: `Close the browser and cleanup resources.

Examples:
  w3pilot browser quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), browserQuitTimeout)
		defer cancel()

		if err := quitBrowser(ctx); err != nil {
			return err
		}

		fmt.Println("Browser closed")
		return nil
	},
}

func init() {
	browserCmd.AddCommand(browserQuitCmd)
	browserQuitCmd.Flags().DurationVar(&browserQuitTimeout, "timeout", 10*time.Second, "Quit timeout")
}
