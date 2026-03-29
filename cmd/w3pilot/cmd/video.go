package cmd

import (
	"github.com/spf13/cobra"
)

// videoCmd represents the video command group
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Video capture commands",
	Long: `Commands for recording video of browser activity.

Examples:
  w3pilot video start                          # Start video recording
  w3pilot video start --dir ./videos           # Start with custom dir
  w3pilot video stop                           # Stop recording`,
}

func init() {
	rootCmd.AddCommand(videoCmd)
}
