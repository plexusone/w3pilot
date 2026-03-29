package cmd

import (
	"github.com/spf13/cobra"
)

// dialogCmd represents the dialog command group
var dialogCmd = &cobra.Command{
	Use:   "dialog",
	Short: "Dialog handling commands",
	Long: `Commands for handling browser dialogs (alert, confirm, prompt).

Examples:
  w3pilot dialog handle --accept               # Accept dialog
  w3pilot dialog handle --dismiss              # Dismiss dialog
  w3pilot dialog handle --text "response"      # Enter prompt text
  w3pilot dialog get                           # Get current dialog info`,
}

func init() {
	rootCmd.AddCommand(dialogCmd)
}
