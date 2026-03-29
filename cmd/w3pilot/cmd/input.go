package cmd

import (
	"github.com/spf13/cobra"
)

// inputCmd represents the input command group
var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Keyboard, mouse, and touch input commands",
	Long: `Commands for low-level keyboard, mouse, and touch input.

Examples:
  w3pilot input key-press Enter                # Press Enter key
  w3pilot input key-type "Hello World"         # Type text via keyboard
  w3pilot input mouse-click 100 200            # Click at coordinates
  w3pilot input mouse-move 100 200             # Move mouse to position
  w3pilot input touch-tap 150 300              # Touch tap at position`,
}

func init() {
	rootCmd.AddCommand(inputCmd)
}
