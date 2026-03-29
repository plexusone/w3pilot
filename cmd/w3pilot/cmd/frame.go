package cmd

import (
	"github.com/spf13/cobra"
)

// frameCmd represents the frame command group
var frameCmd = &cobra.Command{
	Use:   "frame",
	Short: "Frame navigation commands",
	Long: `Commands for working with iframes.

Examples:
  w3pilot frame select "iframe-name"           # Select frame by name
  w3pilot frame select "**/content"            # Select frame by URL
  w3pilot frame main                           # Return to main frame
  w3pilot frame list                           # List all frames`,
}

func init() {
	rootCmd.AddCommand(frameCmd)
}
