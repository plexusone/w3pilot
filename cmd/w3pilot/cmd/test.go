package cmd

import (
	"github.com/spf13/cobra"
)

// testCmd represents the test command group
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing and assertion commands",
	Long: `Commands for testing and verifying page state.

Examples:
  w3pilot test assert-text "Welcome"           # Assert text exists
  w3pilot test assert-element "#login"         # Assert element exists
  w3pilot test assert-url "**/dashboard"       # Assert URL matches
  w3pilot test verify-visible "#modal"         # Verify element visible
  w3pilot test verify-value "#input" "hello"   # Verify input value
  w3pilot test report                          # Generate test report`,
}

func init() {
	rootCmd.AddCommand(testCmd)
}
