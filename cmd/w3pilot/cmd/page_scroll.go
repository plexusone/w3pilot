package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/spf13/cobra"
)

var (
	pageScrollTimeout  time.Duration
	pageScrollSelector string
)

var pageScrollCmd = &cobra.Command{
	Use:   "scroll <direction> [amount]",
	Short: "Scroll the page",
	Long: `Scroll the page or a specific element.

Direction: up, down, left, right
Amount: pixels to scroll (default: full page)

Examples:
  w3pilot page scroll down 500               # Scroll down 500px
  w3pilot page scroll up                     # Scroll up full page
  w3pilot page scroll down --selector "#list" # Scroll element`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		direction := args[0]

		// Validate direction
		switch direction {
		case "up", "down", "left", "right":
			// valid
		default:
			return fmt.Errorf("invalid direction: %s (use up, down, left, right)", direction)
		}

		amount := 0
		if len(args) > 1 {
			var err error
			amount, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid amount: %s", args[1])
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), pageScrollTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		var opts *w3pilot.ScrollOptions
		if pageScrollSelector != "" {
			opts = &w3pilot.ScrollOptions{Selector: pageScrollSelector}
		}

		if err := pilot.Scroll(ctx, direction, amount, opts); err != nil {
			return fmt.Errorf("scroll failed: %w", err)
		}

		if amount > 0 {
			fmt.Printf("Scrolled %s %dpx\n", direction, amount)
		} else {
			fmt.Printf("Scrolled %s\n", direction)
		}
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageScrollCmd)
	pageScrollCmd.Flags().DurationVar(&pageScrollTimeout, "timeout", 10*time.Second, "Timeout")
	pageScrollCmd.Flags().StringVar(&pageScrollSelector, "selector", "", "Scroll within specific element")
}
