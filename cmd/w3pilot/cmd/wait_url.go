package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var waitURLTimeout time.Duration

var waitURLCmd = &cobra.Command{
	Use:   "url <pattern>",
	Short: "Wait for URL to match pattern",
	Long: `Wait for the page URL to match the specified pattern.

Pattern can be a glob or regex.

Examples:
  w3pilot wait url "**/dashboard"
  w3pilot wait url "**/success*"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), waitURLTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.WaitForURL(ctx, pattern, waitURLTimeout); err != nil {
			return fmt.Errorf("wait for URL failed: %w", err)
		}

		url, _ := pilot.URL(ctx)
		fmt.Printf("URL matched: %s\n", url)
		return nil
	},
}

func init() {
	waitCmd.AddCommand(waitURLCmd)
	waitURLCmd.Flags().DurationVar(&waitURLTimeout, "timeout", 30*time.Second, "Timeout")
}
