package cmd

import (
	"github.com/spf13/cobra"
)

// elementCmd represents the element command group
var elementCmd = &cobra.Command{
	Use:   "element",
	Short: "Element interaction commands",
	Long: `Commands for interacting with page elements.

Examples:
  w3pilot element click "#submit"              # Click element
  w3pilot element type "#email" "test@test.com" # Type into input
  w3pilot element fill "#password" "secret"     # Fill input field
  w3pilot element text "#header"               # Get element text
  w3pilot element hover "#menu"                # Hover over element
  w3pilot element check "#checkbox"            # Check checkbox
  w3pilot element visible "#modal"             # Check if visible`,
}

func init() {
	rootCmd.AddCommand(elementCmd)
}
