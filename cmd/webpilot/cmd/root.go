package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	sessionFile string
	verbose     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "webpilot",
	Short: "Browser automation CLI",
	Long: `WebPilot is a browser automation tool that provides:

  - MCP server for AI-assisted browser automation
  - CLI commands for scripted browser control
  - YAML script runner for batch operations

Examples:
  # Start MCP server
  webpilot mcp --headless

  # Launch browser and run commands
  webpilot launch --headless
  webpilot go https://example.com
  webpilot click "#submit"
  webpilot screenshot output.png
  webpilot quit

  # Run a script file
  webpilot run test.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sessionFile, "session", "", "Session file path (default: ~/.webpilot/session.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

// getSessionPath returns the session file path
func getSessionPath() string {
	if sessionFile != "" {
		return sessionFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".webpilot-session.json"
	}
	return fmt.Sprintf("%s/.webpilot/session.json", home)
}
