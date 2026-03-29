package cmd

import (
	"github.com/spf13/cobra"
)

// waitCmd represents the wait command group
var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait condition commands",
	Long: `Commands for waiting on various conditions.

Examples:
  w3pilot wait selector "#modal"               # Wait for element
  w3pilot wait text "Loading complete"         # Wait for text
  w3pilot wait url "**/success"                # Wait for URL pattern
  w3pilot wait load                            # Wait for page load
  w3pilot wait function "window.ready"         # Wait for JS condition`,
}

func init() {
	rootCmd.AddCommand(waitCmd)
}
