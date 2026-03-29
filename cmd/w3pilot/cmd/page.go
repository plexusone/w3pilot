package cmd

import (
	"github.com/spf13/cobra"
)

// pageCmd represents the page command group
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Page navigation and management commands",
	Long: `Commands for navigating, capturing, and managing browser pages.

Examples:
  w3pilot page navigate https://example.com    # Navigate to URL
  w3pilot page back                            # Go back in history
  w3pilot page forward                         # Go forward in history
  w3pilot page reload                          # Reload current page
  w3pilot page screenshot output.png           # Take screenshot
  w3pilot page title                           # Get page title
  w3pilot page url                             # Get current URL
  w3pilot page content                         # Get page HTML content
  w3pilot page scroll down 500                 # Scroll down 500px`,
}

func init() {
	rootCmd.AddCommand(pageCmd)
}
