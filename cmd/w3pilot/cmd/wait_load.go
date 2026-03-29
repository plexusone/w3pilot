package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	waitLoadTimeout time.Duration
	waitLoadState   string
)

var waitLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Wait for page to load",
	Long: `Wait for the page to reach a specific load state.

States: load, domcontentloaded, networkidle

Examples:
  w3pilot wait load
  w3pilot wait load --state networkidle`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), waitLoadTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.WaitForLoad(ctx, waitLoadState, waitLoadTimeout); err != nil {
			return fmt.Errorf("wait for load failed: %w", err)
		}

		fmt.Printf("Page loaded (state: %s)\n", waitLoadState)
		return nil
	},
}

func init() {
	waitCmd.AddCommand(waitLoadCmd)
	waitLoadCmd.Flags().DurationVar(&waitLoadTimeout, "timeout", 30*time.Second, "Timeout")
	waitLoadCmd.Flags().StringVar(&waitLoadState, "state", "load", "Load state to wait for (load, domcontentloaded, networkidle)")
}
