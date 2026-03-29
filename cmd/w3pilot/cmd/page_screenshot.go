//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	pageScreenshotTimeout  time.Duration
	pageScreenshotSelector string
)

var pageScreenshotCmd = &cobra.Command{
	Use:   "screenshot <filename>",
	Short: "Take a screenshot",
	Long: `Capture a screenshot of the current page or a specific element.

Examples:
  w3pilot page screenshot page.png
  w3pilot page screenshot element.png --selector "#main"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), pageScreenshotTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		var data []byte
		var err error

		if pageScreenshotSelector != "" {
			// Element screenshot
			el, findErr := pilot.Find(ctx, pageScreenshotSelector, nil)
			if findErr != nil {
				return fmt.Errorf("element not found: %w", findErr)
			}
			data, err = el.Screenshot(ctx)
		} else {
			// Page screenshot
			data, err = pilot.Screenshot(ctx)
		}

		if err != nil {
			return fmt.Errorf("screenshot failed: %w", err)
		}

		if err := os.WriteFile(filename, data, 0600); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("Screenshot saved: %s (%d bytes)\n", filename, len(data))
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageScreenshotCmd)
	pageScreenshotCmd.Flags().DurationVar(&pageScreenshotTimeout, "timeout", 30*time.Second, "Screenshot timeout")
	pageScreenshotCmd.Flags().StringVar(&pageScreenshotSelector, "selector", "", "Capture specific element instead of page")
}
