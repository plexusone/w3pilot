package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	waitTextTimeout  time.Duration
	waitTextSelector string
)

var waitTextCmd = &cobra.Command{
	Use:   "text <text>",
	Short: "Wait for text to appear",
	Long: `Wait for specific text to appear on the page.

Examples:
  w3pilot wait text "Loading complete"
  w3pilot wait text "Success" --selector "#status"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), waitTextTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		// Build JavaScript to wait for text
		var script string
		if waitTextSelector != "" {
			script = fmt.Sprintf(`
				(function() {
					const el = document.querySelector(%q);
					return el && el.textContent.includes(%q);
				})()
			`, waitTextSelector, text)
		} else {
			script = fmt.Sprintf(`document.body.textContent.includes(%q)`, text)
		}

		if err := pilot.WaitForFunction(ctx, script, waitTextTimeout); err != nil {
			return fmt.Errorf("text not found: %w", err)
		}

		fmt.Printf("Text found: %s\n", text)
		return nil
	},
}

func init() {
	waitCmd.AddCommand(waitTextCmd)
	waitTextCmd.Flags().DurationVar(&waitTextTimeout, "timeout", 30*time.Second, "Timeout")
	waitTextCmd.Flags().StringVar(&waitTextSelector, "selector", "", "Limit search to element matching selector")
}
