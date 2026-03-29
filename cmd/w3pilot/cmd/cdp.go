package cmd

import (
	"github.com/spf13/cobra"
)

// cdpCmd represents the Chrome DevTools Protocol command group
var cdpCmd = &cobra.Command{
	Use:   "cdp",
	Short: "Chrome DevTools Protocol commands",
	Long: `Commands for advanced Chrome DevTools Protocol features.

Examples:
  w3pilot cdp heap snapshot.heapsnapshot       # Take heap snapshot
  w3pilot cdp lighthouse                       # Run Lighthouse audit
  w3pilot cdp coverage start                   # Start code coverage
  w3pilot cdp coverage stop                    # Stop and report coverage
  w3pilot cdp emulate network slow3g           # Emulate slow network
  w3pilot cdp extensions list                  # List extensions`,
}

func init() {
	rootCmd.AddCommand(cdpCmd)
}
