package cmd

import (
	"github.com/spf13/cobra"
)

// jsCmd represents the JavaScript command group
var jsCmd = &cobra.Command{
	Use:   "js",
	Short: "JavaScript execution commands",
	Long: `Commands for executing JavaScript on the page.

Examples:
  w3pilot js eval "document.title"             # Evaluate JavaScript
  w3pilot js add-script "console.log('hi')"    # Add script to page
  w3pilot js add-style "body { color: red }"   # Add CSS to page
  w3pilot js init-script ./setup.js            # Add init script`,
}

func init() {
	rootCmd.AddCommand(jsCmd)
}
