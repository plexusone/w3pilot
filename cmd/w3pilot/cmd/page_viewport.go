package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/spf13/cobra"
)

var (
	pageViewportTimeout time.Duration
	pageViewportWidth   int
	pageViewportHeight  int
)

var pageViewportCmd = &cobra.Command{
	Use:   "viewport",
	Short: "Get or set viewport dimensions",
	Long: `Get or set the browser viewport dimensions.

Examples:
  w3pilot page viewport                      # Get current viewport
  w3pilot page viewport --width 1920 --height 1080  # Set viewport`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageViewportTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		// If width and height are specified, set viewport
		if pageViewportWidth > 0 && pageViewportHeight > 0 {
			viewport := w3pilot.Viewport{
				Width:  pageViewportWidth,
				Height: pageViewportHeight,
			}
			if err := pilot.SetViewport(ctx, viewport); err != nil {
				return fmt.Errorf("failed to set viewport: %w", err)
			}
			fmt.Printf("Viewport set to %dx%d\n", pageViewportWidth, pageViewportHeight)
			return nil
		}

		// Otherwise, get viewport
		viewport, err := pilot.GetViewport(ctx)
		if err != nil {
			return fmt.Errorf("failed to get viewport: %w", err)
		}

		jsonBytes, _ := json.MarshalIndent(viewport, "", "  ")
		fmt.Println(string(jsonBytes))
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageViewportCmd)
	pageViewportCmd.Flags().DurationVar(&pageViewportTimeout, "timeout", 10*time.Second, "Timeout")
	pageViewportCmd.Flags().IntVar(&pageViewportWidth, "width", 0, "Viewport width")
	pageViewportCmd.Flags().IntVar(&pageViewportHeight, "height", 0, "Viewport height")
}
