package cmd

import (
	"github.com/spf13/cobra"
)

// recordCmd represents the record command group
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Action recording commands",
	Long: `Commands for recording and exporting browser actions.

Examples:
  w3pilot record start                         # Start recording
  w3pilot record stop                          # Stop recording
  w3pilot record status                        # Get recording status
  w3pilot record export script.yaml            # Export as script
  w3pilot record clear                         # Clear recorded actions`,
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
