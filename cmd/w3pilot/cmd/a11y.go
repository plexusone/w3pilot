package cmd

import (
	"github.com/spf13/cobra"
)

// a11yCmd represents the accessibility command group
var a11yCmd = &cobra.Command{
	Use:     "a11y",
	Aliases: []string{"accessibility"},
	Short:   "Accessibility commands",
	Long: `Commands for accessibility testing and inspection.

Examples:
  w3pilot a11y snapshot                        # Get accessibility tree
  w3pilot a11y snapshot --interesting-only     # Get interesting nodes only`,
}

func init() {
	rootCmd.AddCommand(a11yCmd)
}
