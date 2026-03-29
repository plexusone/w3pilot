package cmd

import (
	"github.com/spf13/cobra"
)

// traceCmd represents the trace command group
var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Tracing commands",
	Long: `Commands for recording browser traces.

Examples:
  w3pilot trace start                          # Start tracing
  w3pilot trace stop trace.zip                 # Stop and save trace
  w3pilot trace chunk "login flow"             # Add trace chunk
  w3pilot trace group start "auth"             # Start trace group
  w3pilot trace group stop "auth"              # Stop trace group`,
}

func init() {
	rootCmd.AddCommand(traceCmd)
}
